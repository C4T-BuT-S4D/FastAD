package checkers

import (
	"fmt"
	"time"

	"go.temporal.io/sdk/workflow"

	"github.com/c4t-but-s4d/fastad/internal/models"
	checkerpb "github.com/c4t-but-s4d/fastad/pkg/proto/checker"
	gspb "github.com/c4t-but-s4d/fastad/pkg/proto/data/game_state"
	servicespb "github.com/c4t-but-s4d/fastad/pkg/proto/data/services"
	teamspb "github.com/c4t-but-s4d/fastad/pkg/proto/data/teams"
)

const GetWorkflowName = "GetWorkflow"

type GetWorkflowParameters struct {
	GameState *gspb.GameState
	Team      *teamspb.Team
	Service   *servicespb.Service
}

func GetWorkflowDefinition(ctx workflow.Context, params GetWorkflowParameters) error {
	logger := workflow.GetLogger(ctx)
	logger.Info("starting workflow")

	lao := workflow.LocalActivityOptions{
		ScheduleToCloseTimeout: time.Second * 3,
	}
	laoCtx := workflow.WithLocalActivityOptions(ctx, lao)

	// TODO: fail if the last PUT failed.

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

	// TODO: fail if pick hasn't succeeded.
	if pickFlagResult.Flag == nil {
		logger.Info("no flag picked, skipping get")
		return nil
	}

	service := models.NewServiceFromProto(params.Service)

	getActivityCtx := workflow.WithActivityOptions(ctx, workflow.ActivityOptions{
		ScheduleToCloseTimeout: service.CheckerTimeout(checkerpb.Action_ACTION_GET) + checkerKillDelay*2,
	})

	var verdict *Verdict
	if err := workflow.ExecuteActivity(
		getActivityCtx,
		GetActivityName,
		&GetActivityParameters{
			GameState: params.GameState,
			Team:      params.Team,
			Service:   service,
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
