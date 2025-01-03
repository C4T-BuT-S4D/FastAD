package migrations

import (
	"context"
	"fmt"

	"github.com/uptrace/bun"

	"github.com/c4t-but-s4d/fastad/internal/models"
)

//nolint:gochecknoinits // Migrations should be initialized in init functions.
func init() {
	Migrations.MustRegister(func(ctx context.Context, db *bun.DB) error {
		if _, err := db.
			NewCreateTable().
			Model((*models.CheckerExecution)(nil)).
			IfNotExists().
			WithForeignKeys().
			Exec(ctx); err != nil {
			return fmt.Errorf("creating checkers_executions: %w", err)
		}

		if _, err := db.
			NewCreateTable().
			Model((*models.Flag)(nil)).
			IfNotExists().
			WithForeignKeys().
			Exec(ctx); err != nil {
			return fmt.Errorf("creating flags: %w", err)
		}

		if _, err := db.
			NewCreateTable().
			Model((*models.AttackDataSnapshot)(nil)).
			IfNotExists().
			WithForeignKeys().
			Exec(ctx); err != nil {
			return fmt.Errorf("creating attack_data_snapshots: %w", err)
		}

		return nil
	}, func(context.Context, *bun.DB) error {
		return nil
	})
}
