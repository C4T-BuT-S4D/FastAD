package checkers

import (
	"fmt"
	"time"

	"go.temporal.io/sdk/log"
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
	logger := log.With(
		workflow.GetLogger(ctx),
		"team", params.Team.Name,
		"service", params.Service.Name,
		"action", checkerpb.Action_ACTION_CHECK,
	)

	logger.Debug("starting")

	laoCtx := workflow.WithLocalActivityOptions(ctx, workflow.LocalActivityOptions{
		ScheduleToCloseTimeout: time.Second * 3,
	})

	var verdict *Verdict

	var getExecutionResult *GetLastExecutionActivityResult
	if err := workflow.ExecuteLocalActivity(
		laoCtx,
		GetLastExecutionActivityName,
		&GetLastExecutionActivityParameters{
			Action:  checkerpb.Action_ACTION_PUT,
			Team:    params.Team,
			Service: params.Service,
		},
	).Get(ctx, &getExecutionResult); err != nil {
		logger.Error("running get last execution activity", "error", err)
	}

	if lastPut := getExecutionResult.Execution; lastPut != nil && lastPut.Status != checkerpb.Status_STATUS_UP {
		verdict = &Verdict{
			Action:  checkerpb.Action_ACTION_CHECK,
			Status:  lastPut.Status,
			Public:  "last flag PUT failed",
			Private: fmt.Sprintf("last put execution id: %v", lastPut.ID),
		}
		logger.Debug("last put failed", "verdict", verdict)
	} else {
		service := models.NewServiceFromProto(params.Service)

		checkActivityCtx := workflow.WithActivityOptions(ctx, workflow.ActivityOptions{
			ScheduleToCloseTimeout: service.CheckerTimeout(checkerpb.Action_ACTION_CHECK) + checkerKillDelay*2,
		})

		var checkResult CheckActivityResult
		if err := workflow.ExecuteActivity(
			checkActivityCtx,
			CheckActivityName,
			&CheckActivityParameters{
				GameState: params.GameState,
				Team:      params.Team,
				Service:   service,
			},
		).Get(ctx, &checkResult); err != nil {
			logger.Error("running activity", "error", err)
			checkResult.Verdict = &Verdict{
				Action:  checkerpb.Action_ACTION_CHECK,
				Status:  checkerpb.Status_STATUS_CHECK_FAILED,
				Public:  "checker error",
				Private: fmt.Sprintf("running activity: %v", err),
			}
		}

		verdict = checkResult.Verdict
		logger.Debug("checker finished, saving verdict", "verdict", verdict)
	}

	if err := workflow.ExecuteLocalActivity(
		laoCtx,
		SaveVerdictActivityName,
		&SaveVerdictActivityParameters{
			Team:    params.Team,
			Service: params.Service,
			Verdict: verdict,
		},
	).Get(ctx, nil); err != nil {
		logger.Error("running save verdict activity", "error", err)
		return fmt.Errorf("save verdict: %w", err)
	}

	logger.Debug("check finished", "verdict", verdict)

	return nil
}
