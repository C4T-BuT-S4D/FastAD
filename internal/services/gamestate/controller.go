package gamestate

import (
	"context"
	"database/sql"
	"fmt"

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

func (c *Controller) Update(ctx context.Context, req *gspb.UpdateRequest) (*models.GameState, int32, error) {
	gs := &models.GameState{
		ID:                 1,
		StartTime:          req.StartTime.AsTime(),
		TotalRounds:        uint(req.TotalRounds),
		Paused:             req.Paused,
		FlagLifetimeRounds: uint(req.FlagLifetimeRounds),
		RoundDuration:      req.RoundDuration.AsDuration(),
	}

	if req.EndTime != nil {
		gs.EndTime = lo.ToPtr(req.EndTime.AsTime())
	}

	var newVersion int32
	if err := c.db.RunInTx(ctx, &sql.TxOptions{}, func(ctx context.Context, tx bun.Tx) error {
		if err := c.db.
			NewInsert().
			Model(gs).
			On("CONFLICT (id) DO UPDATE").
			Set("start_time = EXCLUDED.start_time").
			Set("end_time = EXCLUDED.end_time").
			Set("total_rounds = EXCLUDED.total_rounds").
			Set("paused = EXCLUDED.paused").
			Set("flag_lifetime_rounds = EXCLUDED.flag_lifetime_rounds").
			Set("round_duration = EXCLUDED.round_duration").
			Returning("*").
			Scan(ctx); err != nil {
			return fmt.Errorf("inserting game state: %w", err)
		}

		var err error
		if newVersion, err = c.Versions.Increment(ctx, tx, VersionKey); err != nil {
			return fmt.Errorf("incrementing version: %w", err)
		}
		return nil
	}); err != nil {
		return nil, 0, fmt.Errorf("in transaction: %w", err)
	}

	return gs, newVersion, nil
}

func (c *Controller) UpdateRound(ctx context.Context, req *gspb.UpdateRoundRequest) (*models.GameState, int32, error) {

	gs := &models.GameState{
		ID:                1,
		RunningRound:      uint(req.RunningRound),
		RunningRoundStart: req.RunningRoundStart.AsTime(),
	}

	var newVersion int32
	if err := c.db.RunInTx(ctx, &sql.TxOptions{}, func(ctx context.Context, tx bun.Tx) error {
		if err := c.db.
			NewInsert().
			Model(gs).
			On("CONFLICT (id) DO UPDATE").
			Set("running_round = EXCLUDED.running_round").
			Set("running_round_start = EXCLUDED.running_round_start").
			Returning("*").
			Scan(ctx); err != nil {
			return fmt.Errorf("updating game state: %w", err)
		}

		var err error
		if newVersion, err = c.Versions.Increment(ctx, tx, VersionKey); err != nil {
			return fmt.Errorf("incrementing version: %w", err)
		}
		return nil
	}); err != nil {
		return nil, 0, fmt.Errorf("in transaction: %w", err)
	}

	return gs, newVersion, nil
}

func (c *Controller) Migrate(ctx context.Context) error {
	if _, err := c.db.
		NewCreateTable().
		IfNotExists().
		Model(&models.GameState{}).
		Exec(ctx); err != nil {
		return fmt.Errorf("creating game state table: %w", err)
	}
	return nil
}
