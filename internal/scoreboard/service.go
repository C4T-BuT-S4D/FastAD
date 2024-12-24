package scoreboard

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/uptrace/bun"
	"go.uber.org/atomic"
	"go.uber.org/zap"

	"github.com/c4t-but-s4d/fastad/internal/models"
	scoreboardpb "github.com/c4t-but-s4d/fastad/pkg/proto/scoreboard"
)

const processorID = "scoreboard"

type Service struct {
	scoreboardpb.UnimplementedScoreboardServiceServer

	db     *bun.DB
	config *Config
	logger *zap.Logger

	state *atomic.Pointer[State]
}

func NewService(db *bun.DB, cfg *Config) *Service {
	return &Service{
		db:     db,
		config: cfg,

		state:  atomic.NewPointer(NewState()),
		logger: zap.L().With(zap.String("component", "scoreboard")),
	}
}

func (s *Service) GetState(context.Context, *scoreboardpb.GetStateRequest) (*scoreboardpb.GetStateResponse, error) {
	return &scoreboardpb.GetStateResponse{Scoreboard: s.state.Load().ToProto()}, nil
}

func (s *Service) Run(ctx context.Context) {
	t := time.NewTicker(s.config.CheckInterval)
	defer t.Stop()

	for {
		select {
		case <-t.C:
			s.logger.Info("checking")
			start := time.Now()
			if err := s.Check(ctx); err != nil {
				s.logger.With(zap.Error(err)).Error("check failed")
			}
			s.logger.With(zap.Duration("duration", time.Since(start))).Info("checked")
		case <-ctx.Done():
			return
		}
	}
}

func (s *Service) Check(ctx context.Context) error {
	stateClone := s.state.Load().Clone()

	if err := s.db.RunInTx(ctx, &sql.TxOptions{}, func(ctx context.Context, tx bun.Tx) error {
		for {
			// TODO: this is a bunch of seq scans, can we do better?
			statesSubquery := tx.
				NewSelect().
				Model(&models.ProcessorState{}).
				Column("entity_id").
				Where("processor_id = ?", processorID)

			executionsIDsSubquery := tx.
				NewSelect().
				Model(&models.CheckerExecution{}).
				Column("id").
				Except(statesSubquery)

			var batch []*models.CheckerExecution
			if err := tx.
				NewSelect().
				Model(&batch).
				Where("id IN (?)", executionsIDsSubquery).
				Order("ce.id ASC").
				Limit(s.config.BatchSize).
				Scan(ctx); err != nil {
				return fmt.Errorf("selecting executions: %w", err)
			}

			if len(batch) == 0 {
				break
			}

			statesToInsert := make([]*models.ProcessorState, 0, len(batch))
			s.logger.With(zap.Int("batch_size", len(batch))).Info("processing executions batch")
			for _, execution := range batch {
				stateClone.Apply(execution)
				statesToInsert = append(statesToInsert, &models.ProcessorState{
					ProcessorID: processorID,
					EntityID:    execution.ID,
					ProcessedAt: time.Now(),
				})
			}

			if _, err := tx.NewInsert().Model(&statesToInsert).Exec(ctx); err != nil {
				return fmt.Errorf("inserting states: %w", err)
			}
		}

		return nil
	}); err != nil {
		return fmt.Errorf("in tx: %w", err)
	}

	s.logger.Info("check done")
	s.logger.Sugar().Debugf("current state: %+v", stateClone)

	s.state.Store(stateClone)

	return nil
}

func (s *Service) RestoreState(ctx context.Context) error {
	start := time.Now()

	if err := s.db.RunInTx(ctx, &sql.TxOptions{}, func(ctx context.Context, tx bun.Tx) error {
		processedExecutionsCount, err := tx.
			NewSelect().
			Model(&models.ProcessorState{}).
			Where("processor_id = ?", processorID).
			Count(ctx)
		if err != nil {
			return fmt.Errorf("counting executions: %w", err)
		}

		s.logger.Info("restoring state from processed executions", zap.Int("executions_count", processedExecutionsCount))

		lastID := -1
		for {
			const batchSize = 1000

			processorStatesSubquery := tx.
				NewSelect().
				Model(&models.ProcessorState{}).
				Column("entity_id").
				Where("processor_id = ? AND entity_id > ?", processorID, lastID).
				Order("entity_id").
				Limit(batchSize)

			var batch []*models.CheckerExecution
			if err := tx.
				NewSelect().
				Model(&models.CheckerExecution{}).
				Where("id IN (?)", processorStatesSubquery).
				Order("id").
				Scan(ctx, &batch); err != nil {
				return fmt.Errorf("fetching executions batch: %w", err)
			}

			if len(batch) == 0 {
				break
			}
			lastID = batch[len(batch)-1].ID

			s.logger.Info("applying batch of executions", zap.Int("batch_size", len(batch)))
			for _, execution := range batch {
				s.state.Load().Apply(execution)
			}

			if len(batch) < batchSize {
				break
			}
		}

		return nil
	}); err != nil {
		return fmt.Errorf("in tx: %w", err)
	}

	s.logger.Info("state restored", zap.Duration("duration", time.Since(start)))

	return nil
}
