package scheduler

import (
	"context"
	safeRand "crypto/rand"
	"database/sql"
	"errors"
	"fmt"
	"io"
	"math/rand/v2"
	"time"

	"github.com/samber/lo"
	"github.com/uptrace/bun"
	"go.temporal.io/api/enums/v1"
	"go.temporal.io/sdk/client"
	"go.uber.org/atomic"
	"go.uber.org/zap"

	"github.com/c4t-but-s4d/fastad/internal/checkers"
	"github.com/c4t-but-s4d/fastad/internal/models"
	checkerpb "github.com/c4t-but-s4d/fastad/pkg/proto/checker"
	gspb "github.com/c4t-but-s4d/fastad/pkg/proto/data/game_state"
	servicespb "github.com/c4t-but-s4d/fastad/pkg/proto/data/services"
	teamspb "github.com/c4t-but-s4d/fastad/pkg/proto/data/teams"
)

const CheckSchedulerInterval = 1 * time.Second

type CheckScheduler struct {
	Team    *teamspb.Team
	Service *servicespb.Service
	Action  checkerpb.Action // Either CHECK or GET.

	gameState *atomic.Pointer[gspb.GameState]

	schedulerID  string
	workflowID   string
	workflowName string

	db             *bun.DB
	temporalClient client.Client
	rnd            *rand.Rand

	logger *zap.Logger
}

func NewCheckScheduler(
	team *teamspb.Team,
	service *servicespb.Service,
	action checkerpb.Action,
	gameState *atomic.Pointer[gspb.GameState],
	db *bun.DB,
	temporalClient client.Client,
) (*CheckScheduler, error) {
	s := &CheckScheduler{
		Team:    team,
		Service: service,
		Action:  action,

		gameState: gameState,

		schedulerID: fmt.Sprintf("check_scheduler_%s_%d_%d", action, team.Id, service.Id),
		workflowID:  fmt.Sprintf("check_workflow_%s_%d_%d", action, team.Id, service.Id),

		db:             db,
		temporalClient: temporalClient,

		logger: zap.L().With(
			zap.String("component", "check_scheduler"),
			zap.Int64("team_id", team.Id),
			zap.Int64("service_id", service.Id),
			zap.String("action", action.String()),
		),
	}

	switch action {
	case checkerpb.Action_ACTION_CHECK:
		s.workflowName = checkers.CheckWorkflowName
	case checkerpb.Action_ACTION_GET:
		s.workflowName = checkers.GetWorkflowName
	default:
		return nil, fmt.Errorf("unknown action: %v", action)
	}

	var seed [32]byte
	if _, err := io.ReadFull(safeRand.Reader, seed[:]); err != nil {
		return nil, fmt.Errorf("reading random seed: %w", err)
	}
	//nolint:gosec // Non-secure generator is OK here.
	s.rnd = rand.New(rand.NewChaCha8(seed))

	return s, nil
}

func (s *CheckScheduler) Run(ctx context.Context) {
	t := time.NewTicker(CheckSchedulerInterval)
	defer t.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-t.C:
			if s.gameState.Load().GetPaused() {
				s.logger.Debug("game is paused, skipping check")
				if err := s.updateStateOnPause(ctx); err != nil {
					s.logger.Error("updating state on pause", zap.Error(err))
				}
				continue
			}

			if err := s.TryRunCheck(ctx); err != nil {
				s.logger.Error("running check", zap.Error(err))
			}
		}
	}
}

func (s *CheckScheduler) TryRunCheck(ctx context.Context) error {
	now := time.Now()

	var state models.SchedulerState
	err := s.db.
		NewSelect().
		Model(&state).
		Where("scheduler_id = ?", s.schedulerID).
		Scan(ctx)
	if errors.Is(err, sql.ErrNoRows) {
		s.logger.Debug("no scheduler state found, initializing")
		state = models.SchedulerState{
			SchedulerID:     s.schedulerID,
			ExpectedNextRun: now,
		}
		if _, err := s.db.NewInsert().Model(&state).Exec(ctx); err != nil {
			return fmt.Errorf("inserting initial scheduler state: %w", err)
		}
	} else if err != nil {
		return fmt.Errorf("getting scheduler state: %w", err)
	}

	checkerTimeout := s.checkerTimeout()

	lowerInterval := checkerTimeout
	upperInterval := 2 * checkerTimeout

	if now.Before(state.ExpectedNextRun) {
		if state.ExpectedNextRun.Sub(now) > upperInterval {
			s.logger.Debug(
				"check is not ready yet, but expected start is too far away",
				zap.Time("expected_next_run", state.ExpectedNextRun),
				zap.Duration("checker_timeout", checkerTimeout),
			)
		} else {
			s.logger.Debug(
				"check is not ready yet",
				zap.Time("expected_next_run", state.ExpectedNextRun),
				zap.Duration("checker_timeout", checkerTimeout),
			)
			return nil
		}
	}

	s.logger.Info("running workflow")
	workflowRun, err := s.temporalClient.ExecuteWorkflow(
		ctx,
		client.StartWorkflowOptions{
			// TODO: split the queues for rounds and checks in production.
			TaskQueue:                                "checkers",
			ID:                                       s.workflowID,
			WorkflowIDConflictPolicy:                 enums.WORKFLOW_ID_CONFLICT_POLICY_FAIL,
			WorkflowIDReusePolicy:                    enums.WORKFLOW_ID_REUSE_POLICY_ALLOW_DUPLICATE,
			WorkflowExecutionErrorWhenAlreadyStarted: true,
		},
		s.workflowName,
		s.workflowParams(),
	)
	if err != nil {
		return fmt.Errorf("executing workflow: %w", err)
	}

	s.logger.Info("workflow started, waiting", zap.String("workflow_id", workflowRun.GetID()))

	if err := workflowRun.Get(ctx, nil); err != nil {
		return fmt.Errorf("waiting for workflow: %w", err)
	}

	delay := s.randomDelay(lowerInterval, upperInterval)
	nextRun := state.ExpectedNextRun.Add(delay)
	if nextRun.Sub(now) > upperInterval {
		nextRun = now.Add(delay)
	}
	s.logger.Info("workflow completed, updating scheduler state", zap.Time("next_run", nextRun))
	if _, err := s.db.
		NewUpdate().
		Model(&state).
		WherePK().
		Set("expected_next_run = ?", nextRun).
		Exec(ctx); err != nil {
		return fmt.Errorf("updating scheduler state: %w", err)
	}

	return nil
}

func (s *CheckScheduler) updateStateOnPause(ctx context.Context) error {
	if _, err := s.db.
		NewInsert().
		Model(&models.SchedulerState{
			SchedulerID:     s.schedulerID,
			ExpectedNextRun: time.Now().Add(CheckSchedulerInterval),
		}).
		On("CONFLICT (scheduler_id) DO UPDATE").
		Set("expected_next_run = excluded.expected_next_run").
		Exec(ctx); err != nil {
		return fmt.Errorf("inserting scheduler state: %w", err)
	}
	return nil
}

func (s *CheckScheduler) randomDelay(minDelay, maxDelay time.Duration) time.Duration {
	minN := minDelay.Nanoseconds()
	maxN := maxDelay.Nanoseconds()
	return time.Duration(s.rnd.Int64N(maxN-minN) + minN)
}

func (s *CheckScheduler) checkerTimeout() time.Duration {
	checkerTimeout := s.Service.Checker.DefaultTimeout.AsDuration()
	if s.Action == checkerpb.Action_ACTION_CHECK {
		checkAction, ok := lo.Find(s.Service.Checker.Actions, func(action *servicespb.Service_Checker_Action) bool {
			return action.Action == s.Action
		})
		if ok && checkAction.Timeout.AsDuration() != 0 {
			checkerTimeout = checkAction.Timeout.AsDuration()
		}
	}
	if s.Action == checkerpb.Action_ACTION_GET {
		getAction, ok := lo.Find(s.Service.Checker.Actions, func(action *servicespb.Service_Checker_Action) bool {
			return action.Action == s.Action
		})
		if ok && getAction.Timeout.AsDuration() != 0 {
			checkerTimeout = getAction.Timeout.AsDuration()
		}
	}
	return checkerTimeout
}

func (s *CheckScheduler) workflowParams() any {
	switch s.Action {
	case checkerpb.Action_ACTION_CHECK:
		return checkers.CheckWorkflowParameters{
			GameState: s.gameState.Load(),
			Team:      s.Team,
			Service:   s.Service,
		}
	case checkerpb.Action_ACTION_GET:
		return checkers.GetWorkflowParameters{
			GameState: s.gameState.Load(),
			Team:      s.Team,
			Service:   s.Service,
		}
	default:
		panic(fmt.Errorf("unknown action: %v", s.Action))
	}
}
