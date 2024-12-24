package checkers

import (
	"fmt"
	"time"

	"go.temporal.io/sdk/workflow"

	"github.com/c4t-but-s4d/fastad/internal/models"
	checkerpb "github.com/c4t-but-s4d/fastad/pkg/proto/checker"
)

const GetWorkflowName = "PutWorkflow"

type GetWorkflowParameters struct {
	GameState *models.GameState
	Team      *models.Team
	Service   *models.Service
}

func GetWorkflowDefinition(ctx workflow.Context, params GetWorkflowParameters) error {
	logger := workflow.GetLogger(ctx)
	logger.Info("starting workflow")

	lao := workflow.LocalActivityOptions{
		ScheduleToCloseTimeout: time.Second * 3,
	}
	laoCtx := workflow.WithLocalActivityOptions(ctx, lao)

	var pickFlagResult *PickGetFlagActivityResult
	if err := workflow.ExecuteLocalActivity(
		laoCtx,
		PickGetFlagActivityName,
		&PickGetFlagActivityParameters{
			GameState: params.GameState,
			Team:      params.Team,
			Service:   params.Service,
		},
	).Get(ctx, &pickFlagResult); err != nil {
		logger.Error("running pick flag activity", "error", err)
	}

	if pickFlagResult.Flag == nil {
		logger.Info("no flag picked, skipping get")
		return nil
	}

	var verdict *Verdict
	if err := workflow.ExecuteActivity(
		ctx,
		GetActivityName,
		&GetActivityParameters{
			GameState: params.GameState,
			Team:      params.Team,
			Service:   params.Service,
			Flag:      pickFlagResult.Flag,
		},
	).Get(ctx, &verdict); err != nil {
		logger.Error("running activity", "error", err)
		verdict = &Verdict{
			Action:  checkerpb.Action_ACTION_GET,
			Status:  checkerpb.Status_STATUS_CHECK_FAILED,
			Public:  "checker error",
			Private: fmt.Sprintf("running activity: %v", err),
		}
	}

	// TODO: save.

	return nil
}
