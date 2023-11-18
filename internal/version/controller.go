package version

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

func (c *Controller) Get(ctx context.Context, name string) (int32, error) {
	var v models.Version
	err := c.db.
		NewSelect().
		Model(&v).
		Where("name = ?", name).
		Scan(ctx)
	if errors.Is(err, sql.ErrNoRows) {
		return 0, nil
	}
	if err != nil {
		return 0, fmt.Errorf("getting version: %w", err)
	}
	return v.Version, nil
}

func (c *Controller) Increment(ctx context.Context, tx bun.Tx, name string) (int32, error) {
	v := &models.Version{
		Name:    name,
		Version: 1,
	}
	if _, err := tx.
		NewInsert().
		Model(&models.Version{
			Name:    name,
			Version: 1,
		}).
		On("CONFLICT (name) DO UPDATE").
		Set("version = v.version + 1").
		Returning("*").
		Exec(ctx); err != nil {
		return 0, fmt.Errorf("inserting version: %w", err)
	}
	return v.Version, nil
}

func (c *Controller) Migrate(ctx context.Context) error {
	if _, err := c.db.
		NewCreateTable().
		IfNotExists().
		Model(&models.Version{}).
		Exec(ctx); err != nil {
		return fmt.Errorf("creating version table: %w", err)
	}
	return nil
}
