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

const GetWorkflowName = "GetWorkflow"

type GetWorkflowParameters struct {
	GameState *gspb.GameState
	Team      *teamspb.Team
	Service   *servicespb.Service
}

func GetWorkflowDefinition(ctx workflow.Context, params GetWorkflowParameters) error {
	logger := log.With(
		workflow.GetLogger(ctx),
		"team", params.Team.Name,
		"service", params.Service.Name,
		"action", checkerpb.Action_ACTION_GET,
	)

	logger.Debug("starting workflow")

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

	logger.Debug("picked flag", "flag", pickFlagResult.Flag)

	// TODO: fail if pick hasn't succeeded.
	if pickFlagResult.Flag == nil {
		logger.Info("no flag picked, skipping get")
		return nil
	}

	service := models.NewServiceFromProto(params.Service)

	getActivityCtx := workflow.WithActivityOptions(ctx, workflow.ActivityOptions{
		ScheduleToCloseTimeout: service.CheckerTimeout(checkerpb.Action_ACTION_GET) + checkerKillDelay*2,
	})

	var getResult GetActivityResult
	if err := workflow.ExecuteActivity(
		getActivityCtx,
		GetActivityName,
		&GetActivityParameters{
			GameState: params.GameState,
			Team:      params.Team,
			Service:   service,
			Flag:      pickFlagResult.Flag,
		},
	).Get(ctx, &getResult); err != nil {
		logger.Error("running activity", "error", err)
		getResult.Verdict = &Verdict{
			Action:  checkerpb.Action_ACTION_GET,
			Status:  checkerpb.Status_STATUS_CHECK_FAILED,
			Public:  "checker error",
			Private: fmt.Sprintf("running activity: %v", err),
		}
	}

	logger.Debug("checker finished, saving verdict", "verdict", getResult.Verdict)

	if err := workflow.ExecuteLocalActivity(
		laoCtx,
		SaveVerdictActivityName,
		&SaveVerdictActivityParameters{
			Team:    params.Team,
			Service: params.Service,
			Verdict: getResult.Verdict,
		},
	).Get(ctx, nil); err != nil {
		logger.Error("running save verdict activity", "error", err)
		return fmt.Errorf("save verdict: %w", err)
	}

	logger.Debug("get finished")

	return nil
}
