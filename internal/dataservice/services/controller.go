package services

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/uptrace/bun"

	"github.com/c4t-but-s4d/fastad/internal/models"
	"github.com/c4t-but-s4d/fastad/internal/version"
)

const VersionKey = "services"

type Controller struct {
	Versions *version.Controller

	db *bun.DB
}

func NewController(db *bun.DB, versionController *version.Controller) *Controller {
	return &Controller{
		Versions: versionController,

		db: db,
	}
}

func (c *Controller) List(ctx context.Context) ([]*models.Service, error) {
	var teams []*models.Service
	if err := c.db.NewSelect().Model(&teams).Scan(ctx); err != nil {
		return nil, fmt.Errorf("getting teams: %w", err)
	}
	return teams, nil
}

func (c *Controller) CreateBatch(ctx context.Context, services []*models.Service) error {
	if len(services) == 0 {
		return nil
	}

	if err := c.db.RunInTx(
		ctx,
		&sql.TxOptions{},
		func(ctx context.Context, tx bun.Tx) error {
			if err := tx.
				NewInsert().
				Model(&services).
				On("CONFLICT (name) DO UPDATE").
				Set("checker_type = EXCLUDED.checker_type").
				Set("checker_path = EXCLUDED.checker_path").
				Set("default_score = EXCLUDED.default_score").
				Set("default_timeout = EXCLUDED.default_timeout").
				Set("actions = EXCLUDED.actions").
				Returning("*").
				Scan(ctx); err != nil {
				return fmt.Errorf("inserting services: %w", err)
			}
			if _, err := c.Versions.Increment(ctx, tx, VersionKey); err != nil {
				return fmt.Errorf("incrementing version: %w", err)
			}
			return nil
		},
	); err != nil {
		return fmt.Errorf("in transaction: %w", err)
	}
	return nil
}