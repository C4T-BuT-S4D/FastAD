package scheduler

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/uptrace/bun"
	"go.temporal.io/api/enums/v1"
	"go.temporal.io/sdk/client"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/c4t-but-s4d/fastad/internal/checkers"
	"github.com/c4t-but-s4d/fastad/internal/models"
	"github.com/c4t-but-s4d/fastad/pkg/clients/gamestate"
	gspb "github.com/c4t-but-s4d/fastad/pkg/proto/data/game_state"
)

const RoundSchedulerID = "round_scheduler"

const (
	RoundLateThreshold       = 5 * time.Second
	CheckRoundInterval       = 1 * time.Second
	RefreshGameStateInterval = 5 * time.Second
)

const RoundWorkflowID = "round_workflow"

type RoundScheduler struct {
	gameStateClient *gamestate.Client
	temporalClient  client.Client

	db *bun.DB

	gameState *gspb.GameState

	logger *zap.Logger
}

func NewRoundScheduler(
	temporalClient client.Client,
	gameStateClient *gamestate.Client,
	db *bun.DB,
	logger *zap.Logger,
) *RoundScheduler {
	return &RoundScheduler{
		temporalClient:  temporalClient,
		gameStateClient: gameStateClient,

		db: db,

		logger: logger.Named("round_scheduler"),
	}
}

func (s *RoundScheduler) Run(ctx context.Context) error {
	if err := s.refreshGameState(ctx); err != nil {
		return fmt.Errorf("fetching initial game state: %w", err)
	}

	refreshTicker := time.NewTicker(RefreshGameStateInterval)
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

			if s.gameState.GetFinished() {
				s.logger.Info("game is finished, skipping round update")
				continue
			}

			if s.gameState.GetPaused() {
				s.logger.Info("game is paused, skipping round update")
				if err := s.skipRun(ctx); err != nil {
					s.logger.Error("skipping scheduler run", zap.Error(err))
				}
				continue
			}

			if s.gameState.GetStartTime().AsTime().After(time.Now()) {
				s.logger.Info("game has not started yet, skipping round update")
				if err := s.skipRun(ctx); err != nil {
					s.logger.Error("skipping scheduler run", zap.Error(err))
				}
				continue
			}

			if s.gameState.GetTotalRounds() > 0 && s.gameState.GetRunningRound() >= s.gameState.GetTotalRounds() {
				s.logger.Info("finishing game, skipping round update")
				if err := s.finishGame(ctx); err != nil {
					s.logger.Error("finishing game", zap.Error(err))
				}
				continue
			}

			if s.gameState.GetEndTime() != nil && time.Now().After(s.gameState.GetEndTime().AsTime()) {
				s.logger.Info("finishing game, skipping round update")
				if err := s.finishGame(ctx); err != nil {
					s.logger.Error("finishing game", zap.Error(err))
				}
				continue
			}

			if err := s.TryRunRound(ctx); err != nil {
				s.logger.Error("running round scheduler", zap.Error(err))
			}

			if err := s.refreshGameState(ctx); err != nil {
				s.logger.Error("refreshing game state", zap.Error(err))
			}
		case <-refreshTicker.C:
			if err := s.refreshGameState(ctx); err != nil {
				s.logger.Error("refreshing game state", zap.Error(err))
				continue
			}
		}
	}
}

func (s *RoundScheduler) TryRunRound(ctx context.Context) error {
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
	s.logger.Debug(
		"checking if ready to run round",
		zap.Duration("round_duration", roundDuration),
		zap.Time("expected_start", state.ExpectedNextRun),
	)

	if now.Before(state.ExpectedNextRun) {
		s.logger.Debug(
			"not ready to run round",
			zap.Time("expected_start", state.ExpectedNextRun),
			zap.Duration("round_duration", roundDuration),
		)
		return nil
	}

	if lag := now.Sub(state.ExpectedNextRun); lag > RoundLateThreshold {
		s.logger.Debug(
			"rounds are running late, trying to catch up",
			zap.Duration("round_duration", roundDuration),
			zap.Time("expected_start", state.ExpectedNextRun),
			zap.Duration("lag", lag),
		)
	}

	s.logger.Info("running round workflow")
	workflowRun, err := s.temporalClient.ExecuteWorkflow(
		ctx,
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

	s.logger.Info("workflow started, waiting for completion", zap.String("run_id", workflowRun.GetRunID()))

	if err := workflowRun.Get(ctx, nil); err != nil {
		return fmt.Errorf("waiting for workflow: %w", err)
	}

	nextRun := state.ExpectedNextRun.Add(roundDuration)
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

func (s *RoundScheduler) skipRun(ctx context.Context) error {
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

func (s *RoundScheduler) refreshGameState(ctx context.Context) error {
	gs, err := s.gameStateClient.Get(ctx)
	if err != nil {
		if status.Code(err) == codes.Unavailable {
			s.gameState = nil
			return nil
		}
		return fmt.Errorf("getting game state: %w", err)
	}
	s.gameState = gs
	return nil
}

func (s *RoundScheduler) finishGame(ctx context.Context) error {
	if _, err := s.gameStateClient.FinishGame(ctx); err != nil {
		return fmt.Errorf("finishing game: %w", err)
	}
	return nil
}
