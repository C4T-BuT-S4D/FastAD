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
		getExecutionResult = &GetLastExecutionActivityResult{}
	}

	if lastPut := getExecutionResult.Execution; lastPut != nil && lastPut.Status != checkerpb.Status_STATUS_UP {
		verdict = &Verdict{
			Action:  checkerpb.Action_ACTION_GET,
			Status:  lastPut.Status,
			Public:  "last flag PUT failed",
			Private: fmt.Sprintf("last put execution id: %v", lastPut.ID),
		}
		logger.Debug("last put failed, failing")
	} else {
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
			pickFlagResult = &PickGetFlagActivityResult{}
		}

		logger.Debug("picked flag", "flag", pickFlagResult.Flag)

		if pickFlagResult.Flag == nil {
			verdict = &Verdict{
				Action:  checkerpb.Action_ACTION_GET,
				Status:  checkerpb.Status_STATUS_CORRUPT,
				Public:  "no flags to get",
				Private: "unable to pick live flag",
			}
			logger.Debug("no flags available, failing")
		} else {
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

			verdict = getResult.Verdict
			logger.Debug("checker finished, saving verdict", "verdict", verdict)
		}
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

	logger.Debug("get finished", "verdict", verdict)

	return nil
}
