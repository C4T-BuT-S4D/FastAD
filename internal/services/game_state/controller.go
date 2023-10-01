package game_state

import (
	"context"
	"fmt"

	"github.com/c4t-but-s4d/fastad/internal/models"
	"github.com/c4t-but-s4d/fastad/internal/version"
	"github.com/uptrace/bun"
)

const VersionKey = "game_state"

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

func (c *Controller) Get(ctx context.Context) (*models.GameState, error) {
	var gs models.GameState
	if err := c.db.NewSelect().Model(&gs).Scan(ctx); err != nil {
		return nil, fmt.Errorf("getting game state: %w", err)
	}
	return &gs, nil
}

func (c *Controller) Migrate(ctx context.Context) error {
	if _, err := c.db.NewCreateTable().IfNotExists().Model(&models.GameState{}).Exec(ctx); err != nil {
		return fmt.Errorf("creating game state table: %w", err)
	}
	return nil
}
