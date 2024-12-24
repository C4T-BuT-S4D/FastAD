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

type Service struct {
	scoreboardpb.UnimplementedScoreboardServiceServer

	db     *bun.DB
	config *Config
	logger *zap.Logger

	state *atomic.Pointer[State]

	lastProcessorIteration time.Time
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
			// We only need to look through the last executions.
			// Some executions can reappear in the "ids holes" if the transaction
			// that inserted them was stuck for some reason.
			// ExecutionTXCreateTimeout is the estimate of this "stuck" time.
			// We select only executions that were created after the last check
			// minus the timeout and minus the threshold ProcessorStateThreshold (just to be safe).
			// This way JOIN with RIGHT IS NULL filter performs better and is
			// simpler than the EXCEPT subquery.

			threshold := s.lastProcessorIteration.
				Add(-s.config.ProcessorStateThreshold).
				Add(-s.config.ExecutionTXCreateTimeout)

			var batch []*models.CheckerExecution
			if err := tx.
				NewSelect().
				Model(&batch).
				Where("created_at > ?", threshold).
				Join("LEFT OUTER JOIN scoreboard_processed_items spi ON spi.checker_execution_id = ce.id").
				Where("spi.id IS NULL").
				Order("ce.id ASC").
				Limit(s.config.BatchSize).
				Scan(ctx); err != nil {
				return fmt.Errorf("selecting executions: %w", err)
			}

			if len(batch) == 0 {
				break
			}

			itemsToInsert := make([]*models.ScoreboardProcessedItem, 0, len(batch))
			s.logger.With(zap.Int("batch_size", len(batch))).Info("processing executions batch")
			for _, execution := range batch {
				stateClone.Apply(execution)
				itemsToInsert = append(itemsToInsert, &models.ScoreboardProcessedItem{
					CheckerExecutionID: execution.ID,
					ProcessedAt:        time.Now(),
				})
			}

			if _, err := tx.NewInsert().Model(&itemsToInsert).Exec(ctx); err != nil {
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
	s.lastProcessorIteration = time.Now()

	return nil
}

func (s *Service) RestoreState(ctx context.Context) error {
	start := time.Now()

	if err := s.db.RunInTx(ctx, &sql.TxOptions{}, func(ctx context.Context, tx bun.Tx) error {
		processedExecutionsCount, err := tx.
			NewSelect().
			Model(&models.ScoreboardProcessedItem{}).
			Count(ctx)
		if err != nil {
			return fmt.Errorf("counting executions: %w", err)
		}

		s.logger.Info("restoring state from processed executions", zap.Int("executions_count", processedExecutionsCount))

		if processedExecutionsCount == 0 {
			return nil
		}

		lastID := -1
		for {
			const batchSize = 1000

			var batch []*models.ScoreboardProcessedItem
			if err := tx.
				NewSelect().
				Model(&models.ScoreboardProcessedItem{}).
				// No need for INNER JOIN as we pretty much know that all processed executions exist.
				Relation("CheckerExecution").
				Where("spi.id > ?", lastID).
				Order("spi.id").
				Limit(batchSize).
				Scan(ctx, &batch); err != nil {
				return fmt.Errorf("fetching processed executions batch: %w", err)
			}

			if len(batch) == 0 {
				break
			}
			lastID = batch[len(batch)-1].ID
			s.lastProcessorIteration = batch[len(batch)-1].ProcessedAt

			s.logger.Info("applying batch of executions", zap.Int("batch_size", len(batch)))
			for _, item := range batch {
				s.state.Load().Apply(item.CheckerExecution)
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
	s.lastProcessorIteration = time.Now()

	return nil
}
