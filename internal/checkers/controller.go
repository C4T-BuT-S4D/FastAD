package checkers

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/uptrace/bun"

	"github.com/c4t-but-s4d/fastad/internal/models"
)

type Controller struct {
	db *bun.DB
}

func NewController(db *bun.DB) *Controller {
	return &Controller{db: db}
}

func (c *Controller) AddFlags(ctx context.Context, flags []*models.Flag) error {
	if _, err := c.db.NewInsert().Model(&flags).Exec(ctx); err != nil {
		return fmt.Errorf("inserting flags: %w", err)
	}
	return nil
}

func (c *Controller) AddCheckerExecutions(ctx context.Context, executions []*models.CheckerExecution) error {
	if _, err := c.db.
		NewInsert().
		Model(&executions).
		On("CONFLICT (execution_id) DO NOTHING").
		Exec(ctx); err != nil {
		return fmt.Errorf("inserting executions: %w", err)
	}
	return nil
}

func (c *Controller) PickFlag(
	ctx context.Context,
	teamID, serviceID int,
	runningRound, lifetimeRounds uint64,
) (*models.Flag, error) {
	minRound := uint64(0)
	if runningRound > lifetimeRounds {
		minRound = runningRound - lifetimeRounds
	}

	var flag models.Flag
	if err := c.db.
		NewSelect().
		Model(&flag).
		Where(
			"team_id = ? AND service_id = ? AND round >= ?",
			teamID,
			serviceID,
			minRound,
		).
		Order("RANDOM()").
		Scan(ctx); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("picking flag: %w", err)
	}

	return &flag, nil
}

func (c *Controller) MigrateDB(ctx context.Context) error {
	if err := c.db.RunInTx(
		ctx,
		&sql.TxOptions{},
		func(ctx context.Context, tx bun.Tx) error {
			if _, err := tx.
				NewCreateTable().
				IfNotExists().
				Model(&models.Flag{}).
				Exec(ctx); err != nil {
				return fmt.Errorf("creating flags table: %w", err)
			}

			if _, err := tx.
				NewCreateTable().
				IfNotExists().
				Model(&models.CheckerExecution{}).
				Exec(ctx); err != nil {
				return fmt.Errorf("creating checker_executions table: %w", err)
			}

			return nil
		},
	); err != nil {
		return fmt.Errorf("migrating: %w", err)
	}

	return nil
}
