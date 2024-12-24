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
		fmt.Print(" [up migration] ")

		if _, err := db.NewCreateTable().
			Model((*models.Version)(nil)).
			IfNotExists().
			Exec(ctx); err != nil {
			return fmt.Errorf("create versions: %w", err)
		}

		if _, err := db.NewCreateTable().
			Model((*models.Team)(nil)).
			IfNotExists().
			Exec(ctx); err != nil {
			return fmt.Errorf("create teams: %w", err)
		}

		if _, err := db.NewCreateTable().
			Model((*models.Service)(nil)).
			IfNotExists().
			Exec(ctx); err != nil {
			return fmt.Errorf("create services: %w", err)
		}

		if _, err := db.NewCreateTable().
			Model((*models.GameState)(nil)).
			IfNotExists().
			Exec(ctx); err != nil {
			return fmt.Errorf("create game_states: %w", err)
		}

		return nil
	}, func(context.Context, *bun.DB) error {
		fmt.Print(" [down migration] ")
		return nil
	})
}
