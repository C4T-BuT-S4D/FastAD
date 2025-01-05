package checkers

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

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

func (c *Controller) SaveAttackDataSnapshot(ctx context.Context, runningRound, lifetimeRounds uint64) error {
	minRound := uint64(0)
	if runningRound > lifetimeRounds {
		minRound = runningRound - lifetimeRounds
	}

	var validFlags []*models.Flag
	if err := c.db.NewSelect().
		Model(&validFlags).
		Where(
			"round >= ? AND f.put_finished is true",
			minRound,
		).
		Order("f.round").
		Scan(ctx); err != nil {
		return fmt.Errorf("getting valid flags: %w", err)
	}

	type teamServiceKey struct {
		teamID, serviceID int
	}
	validFlagsMap := lo.GroupBy(validFlags, func(flag *models.Flag) teamServiceKey {
		return teamServiceKey{flag.TeamID, flag.ServiceID}
	})

	var teams []*models.Team
	if err := c.db.NewSelect().
		Model(&teams).
		Scan(ctx); err != nil {
		return fmt.Errorf("getting teams: %w", err)
	}

	var services []*models.Service
	if err := c.db.NewSelect().
		Model(&services).
		Scan(ctx); err != nil {
		return fmt.Errorf("getting services: %w", err)
	}

	attackDataPayload := make(models.AttackDataPayload)
	for _, service := range services {
		attackDataPayload[service.Name] = make(models.ServiceAttackData)
		for _, team := range teams {
			attackDataPayload[service.Name][team.Address] = []string{}
			key := teamServiceKey{team.ID, service.ID}
			for _, flag := range validFlagsMap[key] {
				attackDataPayload[service.Name][team.Address] = append(
					attackDataPayload[service.Name][team.Address],
					flag.Public,
				)
			}
		}
	}

	snapshot := &models.AttackDataSnapshot{
		CreatedAt: time.Now(),
		Round:     runningRound,
		Payload:   attackDataPayload,
	}
	if _, err := c.db.NewInsert().Model(snapshot).Exec(ctx); err != nil {
		return fmt.Errorf("inserting snapshot: %w", err)
	}

	return nil
}

func (c *Controller) AddCheckerExecutions(ctx context.Context, executions ...*models.CheckerExecution) error {
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
			"team_id = ? AND service_id = ? AND round >= ? AND put_finished is true",
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

func (c *Controller) GetLastExecution(
	ctx context.Context,
	teamID, serviceID int,
	action checkerpb.Action,
) (*models.CheckerExecution, error) {
	var execution models.CheckerExecution
	if err := c.db.
		NewSelect().
		Model(&execution).
		Where(
			"team_id = ? AND service_id = ? AND action = ?",
			teamID,
			serviceID,
			action,
		).
		OrderExpr("created_at DESC").
		Scan(ctx); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			//nolint:nilnil // Easier to handle in the caller.
			return nil, nil
		}
		return nil, fmt.Errorf("getting last put execution: %w", err)
	}

	return &execution, nil
}
