package services

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/uptrace/bun"

	"github.com/c4t-but-s4d/fastad/internal/models"
	"github.com/c4t-but-s4d/fastad/internal/version"
	servicespb "github.com/c4t-but-s4d/fastad/pkg/proto/data/services"
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
	var services []*models.Service
	if err := c.db.NewSelect().Model(&services).Scan(ctx); err != nil {
		return nil, fmt.Errorf("getting services: %w", err)
	}
	return services, nil
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

func (c *Controller) Update(ctx context.Context, req *servicespb.UpdateRequest) (*models.Service, int, error) {
	var service models.Service
	var newVersion int
	if err := c.db.RunInTx(
		ctx,
		&sql.TxOptions{},
		func(ctx context.Context, tx bun.Tx) error {
			query := c.db.NewUpdate().
				Model(&models.Service{}).
				Where("id = ?", req.GetId()).
				Returning("*")

			if req.Disabled != nil {
				query.Set("disabled = ?", req.GetDisabled())
			}

			if err := query.Model(&service).Scan(ctx); err != nil {
				return fmt.Errorf("updating service: %w", err)
			}

			var err error
			if newVersion, err = c.Versions.Increment(ctx, tx, VersionKey); err != nil {
				return fmt.Errorf("incrementing version: %w", err)
			}

			return nil
		},
	); err != nil {
		return nil, 0, fmt.Errorf("in transaction: %w", err)
	}

	return &service, newVersion, nil
}
