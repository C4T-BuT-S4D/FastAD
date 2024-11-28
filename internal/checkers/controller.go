package checkers

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/c4t-but-s4d/fastad/internal/models"
	"github.com/uptrace/bun"
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
