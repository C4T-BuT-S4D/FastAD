package gamestate

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/samber/lo"
	"github.com/uptrace/bun"

	"github.com/c4t-but-s4d/fastad/internal/models"
	"github.com/c4t-but-s4d/fastad/internal/version"
	gspb "github.com/c4t-but-s4d/fastad/pkg/proto/data/game_state"
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

func (c *Controller) Update(ctx context.Context, req *gspb.UpdateRequest) (*models.GameState, int, error) {
	var newVersion int
	var gs *models.GameState
	if err := c.db.RunInTx(ctx, &sql.TxOptions{}, func(ctx context.Context, tx bun.Tx) error {
		var err error
		if gs, newVersion, err = c.updateImpl(ctx, tx, func(query *bun.UpdateQuery) *bun.UpdateQuery {
			var endTime *time.Time
			if req.EndTime != nil {
				endTime = lo.ToPtr(req.EndTime.AsTime())
			}

			return query.
				Set("start_time = ?", req.StartTime.AsTime()).
				Set("end_time = ?", endTime).
				Set("total_rounds = ?", req.TotalRounds).
				Set("paused = ?", req.Paused).
				Set("flag_lifetime_rounds = ?", req.FlagLifetimeRounds).
				Set("round_duration = ?", req.RoundDuration.AsDuration()).
				Set("hardness = ?", req.Hardness).
				Set("inflation = ?", req.Inflation)
		}); err != nil {
			return fmt.Errorf("updating game state: %w", err)
		}
		return nil
	}); err != nil {
		return nil, 0, fmt.Errorf("in transaction: %w", err)
	}

	return gs, newVersion, nil
}

func (c *Controller) UpdateRound(ctx context.Context, req *gspb.UpdateRoundRequest) (*models.GameState, int, error) {
	var newVersion int
	var gs *models.GameState
	if err := c.db.RunInTx(ctx, &sql.TxOptions{}, func(ctx context.Context, tx bun.Tx) error {
		var err error
		if gs, newVersion, err = c.updateImpl(ctx, tx, func(query *bun.UpdateQuery) *bun.UpdateQuery {
			return query.
				Set("running_round = ?", req.RunningRound).
				Set("running_round_start = ?", req.RunningRoundStart.AsTime())
		}); err != nil {
			return fmt.Errorf("updating game state: %w", err)
		}
		return nil
	}); err != nil {
		return nil, 0, fmt.Errorf("in transaction: %w", err)
	}

	return gs, newVersion, nil
}

func (c *Controller) FinishGame(ctx context.Context) (*models.GameState, int, error) {
	var newVersion int
	var gs *models.GameState
	if err := c.db.RunInTx(ctx, &sql.TxOptions{}, func(ctx context.Context, tx bun.Tx) error {
		var err error
		if gs, newVersion, err = c.updateImpl(ctx, tx, func(query *bun.UpdateQuery) *bun.UpdateQuery {
			return query.Set("finished = true")
		}); err != nil {
			return fmt.Errorf("updating game state: %w", err)
		}
		return nil
	}); err != nil {
		return nil, 0, fmt.Errorf("in transaction: %w", err)
	}

	return gs, newVersion, nil
}

func (c *Controller) updateImpl(
	ctx context.Context,
	tx bun.Tx,
	queryMod func(query *bun.UpdateQuery) *bun.UpdateQuery,
) (*models.GameState, int, error) {
	if err := tx.
		NewSelect().
		Model(&models.GameState{}).
		Where("id = 1").
		For("UPDATE").
		Scan(ctx); err != nil {
		return nil, 0, fmt.Errorf("locking game state: %w", err)
	}

	var gs models.GameState
	query := tx.NewUpdate().Model(&gs).Where("id = 1").Returning("*")
	query = queryMod(query)

	if err := query.Scan(ctx); err != nil {
		return nil, 0, fmt.Errorf("updating game state: %w", err)
	}

	newVersion, err := c.Versions.Increment(ctx, tx, VersionKey)
	if err != nil {
		return nil, 0, fmt.Errorf("incrementing version: %w", err)
	}
	return &gs, newVersion, nil
}
