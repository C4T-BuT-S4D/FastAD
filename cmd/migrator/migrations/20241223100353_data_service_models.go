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
		if _, err := db.NewCreateTable().
			Model((*models.Version)(nil)).
			IfNotExists().
			Exec(ctx); err != nil {
			return fmt.Errorf("creating versions: %w", err)
		}

		if _, err := db.NewCreateTable().
			Model((*models.Team)(nil)).
			IfNotExists().
			Exec(ctx); err != nil {
			return fmt.Errorf("creating teams: %w", err)
		}

		if _, err := db.NewCreateTable().
			Model((*models.Service)(nil)).
			IfNotExists().
			Exec(ctx); err != nil {
			return fmt.Errorf("creating services: %w", err)
		}

		if _, err := db.NewCreateTable().
			Model((*models.GameState)(nil)).
			IfNotExists().
			Exec(ctx); err != nil {
			return fmt.Errorf("creating game_states: %w", err)
		}

		return nil
	}, func(context.Context, *bun.DB) error {
		return nil
	})
}
