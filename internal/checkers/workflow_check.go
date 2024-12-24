package checkers

import (
	"fmt"

	"go.temporal.io/sdk/workflow"

	"github.com/c4t-but-s4d/fastad/internal/models"
	checkerpb "github.com/c4t-but-s4d/fastad/pkg/proto/checker"
)

const CheckWorkflowName = "CheckWorkflow"

type CheckWorkflowParameters struct {
	GameState *models.GameState
	Team      *models.Team
	Service   *models.Service
}

func CheckWorkflowDefinition(ctx workflow.Context, params CheckWorkflowParameters) error {
	logger := workflow.GetLogger(ctx)
	logger.Info("starting workflow")

	var verdict *Verdict
	if err := workflow.ExecuteActivity(
		ctx,
		CheckActivityName,
		&CheckActivityParameters{
			GameState: params.GameState,
			Team:      params.Team,
			Service:   params.Service,
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
