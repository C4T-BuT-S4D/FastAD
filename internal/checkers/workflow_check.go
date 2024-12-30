package checkers

import (
	"fmt"

	"go.temporal.io/sdk/workflow"

	"github.com/c4t-but-s4d/fastad/internal/models"
	checkerpb "github.com/c4t-but-s4d/fastad/pkg/proto/checker"
	gspb "github.com/c4t-but-s4d/fastad/pkg/proto/data/game_state"
	servicespb "github.com/c4t-but-s4d/fastad/pkg/proto/data/services"
	teamspb "github.com/c4t-but-s4d/fastad/pkg/proto/data/teams"
)

const CheckWorkflowName = "CheckWorkflow"

type CheckWorkflowParameters struct {
	GameState *gspb.GameState
	Team      *teamspb.Team
	Service   *servicespb.Service
}

func CheckWorkflowDefinition(ctx workflow.Context, params CheckWorkflowParameters) error {
	logger := workflow.GetLogger(ctx)
	logger.Info("starting workflow")

	// TODO: fail if the last PUT failed.
	service := models.NewServiceFromProto(params.Service)

	checkActivityCtx := workflow.WithActivityOptions(ctx, workflow.ActivityOptions{
		ScheduleToCloseTimeout: service.CheckerTimeout(checkerpb.Action_ACTION_CHECK) + checkerKillDelay*2,
	})

	var verdict *Verdict
	if err := workflow.ExecuteActivity(
		checkActivityCtx,
		CheckActivityName,
		&CheckActivityParameters{
			GameState: params.GameState,
			Team:      params.Team,
			Service:   service,
		},
	).Get(ctx, &verdict); err != nil {
		logger.Error("running activity", "error", err)
		verdict = &Verdict{
			Action:  checkerpb.Action_ACTION_CHECK,
			Status:  checkerpb.Status_STATUS_CHECK_FAILED,
			Public:  "checker error",
			Private: fmt.Sprintf("running activity: %v", err),
		}
	}

	// TODO: save.

	return nil
}
