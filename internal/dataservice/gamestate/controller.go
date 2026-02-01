package gamestate

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/uptrace/bun"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

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

func (c *Controller) Create(ctx context.Context, req *gspb.CreateRequest) (*models.GameState, int, error) {
	var newVersion int
	gs := models.NewGameStateFromProto(req.GetGameState())
	if err := c.db.RunInTx(ctx, &sql.TxOptions{}, func(ctx context.Context, tx bun.Tx) error {
		var err error
		if newVersion, err = c.Versions.Increment(ctx, tx, VersionKey); err != nil {
			return fmt.Errorf("incrementing version: %w", err)
		}
		if newVersion > 1 {
			return status.Error(codes.FailedPrecondition, "game state already exists")
		}
		gs.ID = 1
		if err := tx.NewInsert().Model(gs).Returning("*").Scan(ctx); err != nil {
			return fmt.Errorf("inserting game state: %w", err)
		}

		return nil
	}); err != nil {
		return nil, 0, fmt.Errorf("in transaction: %w", err)
	}

	return gs, newVersion, nil
}

func (c *Controller) Update(ctx context.Context, req *gspb.UpdateRequest) (*models.GameState, int, error) {
	var newVersion int
	var gs *models.GameState
	if err := c.db.RunInTx(ctx, &sql.TxOptions{}, func(ctx context.Context, tx bun.Tx) error {
		var err error
		if gs, newVersion, err = c.updateImpl(ctx, tx, func(query *bun.UpdateQuery) *bun.UpdateQuery {
			if req.GetStartTime() != nil {
				query = query.Set("start_time = ?", req.GetStartTime().AsTime())
			}
			if req.GetEndTime() != nil {
				if req.GetEndTime().AsTime().IsZero() {
					query = query.Set("end_time = NULL")
				} else {
					query = query.Set("end_time = ?", req.GetEndTime().AsTime())
				}
			}

			if req.TotalRounds != nil {
				if req.GetTotalRounds() > 0 {
					query = query.Set("total_rounds = ?", req.GetTotalRounds())
				} else {
					query = query.Set("total_rounds = NULL")
				}
			}

			if req.GetStatus() != gspb.GameStatus_GAME_STATUS_UNSPECIFIED {
				query = query.Set("status = ?", req.GetStatus())
			}

			if req.GetFlagLifetimeRounds() > 0 {
				query = query.Set("flag_lifetime_rounds = ?", req.GetFlagLifetimeRounds())
			}

			if req.GetRoundDuration().AsDuration() > 0 {
				query = query.Set("round_duration = ?", req.GetRoundDuration().AsDuration())
			}

			if req.Hardness != nil {
				query = query.Set("hardness = ?", req.GetHardness())
			}

			if req.Inflation != nil {
				query = query.Set("inflation = ?", req.GetInflation())
			}

			return query
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
				Set("running_round = ?", req.GetRunningRound()).
				Set("running_round_start = ?", req.GetRunningRoundStart().AsTime())
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
			return query.Set("status = ?", gspb.GameStatus_GAME_STATUS_FINISHED)
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
