package teams

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/c4t-but-s4d/fastad/internal/models"
	"github.com/c4t-but-s4d/fastad/internal/version"
	"github.com/uptrace/bun"
)

const VersionKey = "teams"

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

func (c *Controller) List(ctx context.Context) ([]*models.Team, error) {
	var teams []*models.Team
	if err := c.db.NewSelect().Model(&teams).Scan(ctx); err != nil {
		return nil, fmt.Errorf("getting teams: %w", err)
	}
	return teams, nil
}

func (c *Controller) CreateBatch(ctx context.Context, teams []*models.Team) error {
	if len(teams) == 0 {
		return nil
	}

	if err := c.db.RunInTx(
		ctx,
		&sql.TxOptions{},
		func(ctx context.Context, tx bun.Tx) error {
			if _, err := c.db.
				NewInsert().
				Model(&teams).
				On("CONFLICT (name) DO UPDATE").
				Set("address = EXCLUDED.address").
				Set("token = EXCLUDED.token").
				Set("labels = EXCLUDED.labels").
				Exec(ctx); err != nil {
				return fmt.Errorf("inserting teams: %w", err)
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

func (c *Controller) Migrate(ctx context.Context) error {
	if _, err := c.db.NewCreateTable().IfNotExists().Model(&models.Team{}).Exec(ctx); err != nil {
		return fmt.Errorf("creating teams table: %w", err)
	}
	return nil
}
