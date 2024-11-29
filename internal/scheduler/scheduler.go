package scheduler

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/uptrace/bun"
	"go.temporal.io/api/enums/v1"
	"go.temporal.io/sdk/client"

	"github.com/c4t-but-s4d/fastad/internal/checkers"
	"github.com/c4t-but-s4d/fastad/internal/clients/gamestate"
	"github.com/c4t-but-s4d/fastad/internal/models"
	gspb "github.com/c4t-but-s4d/fastad/pkg/proto/data/game_state"
)

const RoundSchedulerID = "round_scheduler"

const (
	RoundLateThreshold = 5 * time.Second
	CheckRoundInterval = 1 * time.Second
)

const RoundWorkflowID = "round_workflow"

type Scheduler struct {
	refreshInterval time.Duration

	gameStateClient *gamestate.Client
	temporalClient  client.Client

	db *bun.DB

	gameState *gspb.GameState
	gsVersion int64

	logger *logrus.Entry
}

func New(
	refreshInterval time.Duration,
	temporalClient client.Client,
	gameStateClient *gamestate.Client,
	db *bun.DB,
) *Scheduler {
	return &Scheduler{
		refreshInterval: refreshInterval,
		temporalClient:  temporalClient,
		gameStateClient: gameStateClient,

		db: db,

		logger: logrus.WithField("component", "scheduler"),
	}
}

func (s *Scheduler) Run(ctx context.Context) error {
	if err := s.refreshGameState(ctx); err != nil {
		return fmt.Errorf("fetching initial game state: %w", err)
	}

	refreshTicker := time.NewTicker(s.refreshInterval)
	defer refreshTicker.Stop()

	checkTicker := time.NewTicker(CheckRoundInterval)
	defer checkTicker.Stop()

	for {
		select {
		case <-ctx.Done():
			return nil
		case <-checkTicker.C:
			if s.gameState == nil {
				s.logger.Info("game state is not initialized, skipping round update")
				continue
			}
			if s.gameState.GetPaused() {
				s.logger.Info("game is paused, skipping round update")
				if err := s.updateStateOnPause(ctx); err != nil {
					s.logger.WithError(err).Error("updating scheduler state on pause")
				}
				continue
			}

			if err := s.TryRunRound(ctx); err != nil {
				s.logger.WithError(err).Error("running round scheduler")
			}
		case <-refreshTicker.C:
			if err := s.refreshGameState(ctx); err != nil {
				s.logger.WithError(err).Error("refreshing game state")
				continue
			}
		}
	}
}

func (s *Scheduler) TryRunRound(ctx context.Context) error {
	now := time.Now()

	var state models.SchedulerState
	err := s.db.
		NewSelect().
		Model(&state).
		Where("scheduler_id = ?", RoundSchedulerID).
		Scan(ctx)
	if errors.Is(err, sql.ErrNoRows) {
		s.logger.Debug("no scheduler state found, initializing")
		state = models.SchedulerState{
			SchedulerID:     RoundSchedulerID,
			ExpectedNextRun: now,
		}
		if _, err := s.db.NewInsert().Model(&state).Exec(ctx); err != nil {
			return fmt.Errorf("inserting initial scheduler state: %w", err)
		}
	} else if err != nil {
		return fmt.Errorf("getting scheduler state: %w", err)
	}

	roundDuration := s.gameState.GetRoundDuration().AsDuration()
	s.logger.WithFields(logrus.Fields{
		"round_duration": roundDuration,
		"expected_start": state.ExpectedNextRun,
	}).Debug("checking if ready to run round")

	if now.Before(state.ExpectedNextRun) {
		s.logger.WithFields(logrus.Fields{
			"expected_start": state.ExpectedNextRun,
			"round_duration": roundDuration,
		}).Debug("not ready to run round")
		return nil
	}

	if lag := now.Sub(state.ExpectedNextRun); lag > RoundLateThreshold {
		s.logger.WithFields(logrus.Fields{
			"round_duration": roundDuration,
			"expected_start": state.ExpectedNextRun,
			"lag":            lag,
		}).Warn("rounds are running late, trying to catch up")
	}

	s.logger.Info("running round workflow")
	workflowRun, err := s.temporalClient.ExecuteWorkflow(
		context.Background(),
		client.StartWorkflowOptions{
			TaskQueue:                                "checkers",
			ID:                                       RoundWorkflowID,
			WorkflowIDConflictPolicy:                 enums.WORKFLOW_ID_CONFLICT_POLICY_FAIL,
			WorkflowIDReusePolicy:                    enums.WORKFLOW_ID_REUSE_POLICY_ALLOW_DUPLICATE,
			WorkflowExecutionErrorWhenAlreadyStarted: true,
		},
		checkers.RoundWorkflowName,
		checkers.RoundWorkflowParameters{},
	)
	if err != nil {
		return fmt.Errorf("executing workflow: %w", err)
	}

	s.logger.
		WithField("run_id", workflowRun.GetRunID()).
		Info("workflow started, waiting for completion")

	if err := workflowRun.Get(ctx, nil); err != nil {
		return fmt.Errorf("waiting for workflow: %w", err)
	}

	nextRun := state.ExpectedNextRun.Add(roundDuration)
	s.logger.Infof("workflow completed, updating scheduler state, next run: %v", nextRun)
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

func (s *Scheduler) MigrateDB(ctx context.Context) error {
	if _, err := s.db.
		NewCreateTable().
		IfNotExists().
		Model(&models.SchedulerState{}).
		Exec(ctx); err != nil {
		return fmt.Errorf("creating scheduler state table: %w", err)
	}
	return nil
}

func (s *Scheduler) updateStateOnPause(ctx context.Context) error {
	if _, err := s.db.
		NewInsert().
		Model(&models.SchedulerState{
			SchedulerID:     RoundSchedulerID,
			ExpectedNextRun: time.Now().Add(CheckRoundInterval),
		}).
		On("CONFLICT (scheduler_id) DO UPDATE").
		Set("expected_next_run = excluded.expected_next_run").
		Exec(ctx); err != nil {
		return fmt.Errorf("inserting scheduler state: %w", err)
	}
	return nil
}

func (s *Scheduler) refreshGameState(ctx context.Context) error {
	response, err := s.gameStateClient.RawClient().Get(ctx, &gspb.GetRequest{})
	if err != nil {
		return fmt.Errorf("getting game state: %w", err)
	}

	if s.gameState == nil || s.gsVersion != response.GetVersion().GetVersion() {
		s.gameState = response.GetGameState()
		s.logger.WithFields(logrus.Fields{
			"old_version": s.gsVersion,
			"new_version": response.GetVersion().GetVersion(),
		}).Info("updated game state")
		s.gsVersion = response.GetVersion().GetVersion()
	}

	return nil
}
