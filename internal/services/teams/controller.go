package teams

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/c4t-but-s4d/fastad/internal/models"
	"github.com/redis/go-redis/v9"
	"github.com/uptrace/bun"
)

const lastUpdateRedisKey = "teams::last_update"

type Controller struct {
	db *bun.DB
	r  *redis.Client
}

func NewController(db *bun.DB, r *redis.Client) *Controller {
	return &Controller{db: db, r: r}
}

func (c *Controller) LastUpdate(ctx context.Context) (int64, error) {
	t, err := c.r.Get(ctx, lastUpdateRedisKey).Time()
	if err == nil {
		return t.UnixNano(), nil
	}
	if errors.Is(err, redis.Nil) {
		return 0, nil
	}
	return 0, fmt.Errorf("getting last update time: %w", err)
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
	if _, err := c.db.NewInsert().Model(&teams).Exec(ctx); err != nil {
		return fmt.Errorf("inserting teams: %w", err)
	}
	if err := c.r.Set(ctx, lastUpdateRedisKey, time.Now(), 0).Err(); err != nil {
		return fmt.Errorf("setting last update time: %w", err)
	}
	return nil
}

func (c *Controller) Migrate(ctx context.Context) error {
	if _, err := c.db.NewCreateTable().IfNotExists().Model(&models.Team{}).Exec(ctx); err != nil {
		return fmt.Errorf("creating teams table: %w", err)
	}
	if err := c.r.Set(ctx, lastUpdateRedisKey, time.Now(), 0).Err(); err != nil {
		return fmt.Errorf("setting last update time: %w", err)
	}
	return nil
}
