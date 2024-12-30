package checkers

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/samber/lo"
	"github.com/uptrace/bun"

	"github.com/c4t-but-s4d/fastad/internal/models"
	checkerpb "github.com/c4t-but-s4d/fastad/pkg/proto/checker"
)

type Controller struct {
	db *bun.DB
}

func NewController(db *bun.DB) *Controller {
	return &Controller{db: db}
}

func (c *Controller) AddFlags(ctx context.Context, flags []*models.Flag) error {
	if _, err := c.db.NewInsert().Model(&flags).Exec(ctx); err != nil {
		return fmt.Errorf("inserting flags: %w", err)
	}
	return nil
}

func (c *Controller) SavePutExecutions(ctx context.Context, putResults []*PutActivityResult) error {
	var flagsToRemove []int
	flagsToMark := make([]*models.Flag, 0, len(putResults))
	executions := make([]*models.CheckerExecution, 0, len(putResults))
	for _, putResult := range putResults {
		execution := &models.CheckerExecution{
			ExecutionID: fmt.Sprintf("put-flag-%d", putResult.FlagInfo.Flag.ID),
			TeamID:      putResult.FlagInfo.Team.ID,
			ServiceID:   putResult.FlagInfo.Service.ID,
			Action:      checkerpb.Action_ACTION_PUT,
			Status:      putResult.Verdict.Status,
			Public:      putResult.Verdict.Public,
			Private:     putResult.Verdict.Private,
			Command:     putResult.Verdict.Command,
		}
		executions = append(executions, execution)
		if putResult.Verdict.Status != checkerpb.Status_STATUS_UP {
			flagsToRemove = append(flagsToRemove, putResult.FlagInfo.Flag.ID)
		} else {
			flagsToMark = append(flagsToMark, &models.Flag{
				ID:          putResult.FlagInfo.Flag.ID,
				Public:      putResult.Verdict.Public,
				Private:     putResult.Verdict.Private,
				PutFinished: true,
			})
		}
	}

	if len(executions) == 0 {
		return nil
	}

	if err := c.db.RunInTx(ctx, &sql.TxOptions{}, func(ctx context.Context, tx bun.Tx) error {
		if _, err := tx.
			NewInsert().
			Model(&executions).
			On("CONFLICT (execution_id) DO NOTHING").
			Exec(ctx); err != nil {
			return fmt.Errorf("inserting executions: %w", err)
		}

		if len(flagsToRemove) > 0 {
			if _, err := tx.
				NewDelete().
				Model(&models.Flag{}).
				Where("id IN (?) AND not put_finished", bun.In(flagsToRemove)).
				Exec(ctx); err != nil {
				return fmt.Errorf("deleting flags: %w", err)
			}
		}

		if len(flagsToMark) > 0 {
			for _, chunk := range lo.Chunk(flagsToMark, 200) {
				if _, err := tx.
					NewUpdate().
					Model(&chunk).
					Column("public", "private", "put_finished").
					Bulk().
					Exec(ctx); err != nil {
					return fmt.Errorf("updating flags: %w", err)
				}
			}
		}

		return nil
	}); err != nil {
		return fmt.Errorf("in tx: %w", err)
	}
	return nil
}

func (c *Controller) AddCheckerExecutions(ctx context.Context, executions []*models.CheckerExecution) error {
	if _, err := c.db.
		NewInsert().
		Model(&executions).
		On("CONFLICT (execution_id) DO NOTHING").
		Exec(ctx); err != nil {
		return fmt.Errorf("inserting executions: %w", err)
	}
	return nil
}

func (c *Controller) PickFlag(
	ctx context.Context,
	teamID, serviceID int,
	runningRound, lifetimeRounds uint64,
) (*models.Flag, error) {
	minRound := uint64(0)
	if runningRound > lifetimeRounds {
		minRound = runningRound - lifetimeRounds
	}

	var flag models.Flag
	if err := c.db.
		NewSelect().
		Model(&flag).
		Where(
			"team_id = ? AND service_id = ? AND round >= ?",
			teamID,
			serviceID,
			minRound,
		).
		OrderExpr("RANDOM()").
		Scan(ctx); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			//nolint:nilnil // Easier to handle in the caller.
			return nil, nil
		}
		return nil, fmt.Errorf("picking flag: %w", err)
	}

	return &flag, nil
}
