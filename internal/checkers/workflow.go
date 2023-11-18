package checkers

import (
	"fmt"
	"time"

	"go.temporal.io/sdk/workflow"

	"github.com/c4t-but-s4d/fastad/internal/models"
	checkerpb "github.com/c4t-but-s4d/fastad/pkg/proto/checker"
)

type WorkflowParameters struct {
	GameState *models.GameState
}

func WorkflowDefinition(ctx workflow.Context, params WorkflowParameters) error {
	logger := workflow.GetLogger(ctx)
	logger.Info("starting workflow")

	lao := workflow.LocalActivityOptions{
		ScheduleToCloseTimeout: time.Second * 3,
	}
	laoCtx := workflow.WithLocalActivityOptions(ctx, lao)

	var fetchDataResult *ActivityFetchDataResult
	if err := workflow.ExecuteLocalActivity(
		laoCtx,
		ActivityFetchDataName,
		&ActivityFetchDataParameters{
			GameState: params.GameState,
		},
	).Get(ctx, &fetchDataResult); err != nil {
		return fmt.Errorf("running fetch data activity: %w", err)
	}

	wg := workflow.NewWaitGroup(ctx)
	wg.Add(len(fetchDataResult.Teams) * len(fetchDataResult.Services))
	for _, team := range fetchDataResult.Teams {
		for _, service := range fetchDataResult.Services {
			team := team
			service := service
			workflow.Go(ctx, func(ctx workflow.Context) {
				defer wg.Done()
				runCheckers(ctx, team, service)
			})
		}
	}
	wg.Wait(ctx)

	return nil
}

func runCheckers(ctx workflow.Context, team *models.Team, service *models.Service) {
	logger := workflow.GetLogger(ctx)

	checkActivityCtx := workflow.WithActivityOptions(ctx, workflow.ActivityOptions{
		ScheduleToCloseTimeout: service.CheckerTimeout(checkerpb.Action_ACTION_CHECK) + checkerKillDelay*2,
	})
	putActivityCtx := workflow.WithActivityOptions(ctx, workflow.ActivityOptions{
		ScheduleToCloseTimeout: service.CheckerTimeout(checkerpb.Action_ACTION_PUT) + checkerKillDelay*2,
	})
	getActivityCtx := workflow.WithActivityOptions(ctx, workflow.ActivityOptions{
		ScheduleToCloseTimeout: service.CheckerTimeout(checkerpb.Action_ACTION_GET) + checkerKillDelay*2,
	})

	putCount := service.GetRunCount(checkerpb.Action_ACTION_PUT)
	getCount := service.GetRunCount(checkerpb.Action_ACTION_GET)

	var checkResult *CheckActivityResult
	putResults := make([]*PutActivityResult, 0, putCount)
	getResults := make([]*GetActivityResult, 0, getCount)

	if err := workflow.ExecuteActivity(
		checkActivityCtx,
		ActivityCheckName,
		CheckActivityParameters{
			Team:    team,
			Service: service,
		},
	).Get(ctx, &checkResult); err != nil {
		checkResult = &CheckActivityResult{
			Verdict: &models.CheckerVerdict{
				Action:  checkerpb.Action_ACTION_CHECK,
				Status:  checkerpb.Status_STATUS_CHECK_FAILED,
				Public:  "internal error",
				Private: fmt.Sprintf("check activity err: %v", err),
			},
		}
	}

	if checkResult.Verdict.IsUp() {
		putResultsChan := workflow.NewBufferedChannel(ctx, putCount)

		for i := 0; i < putCount; i++ {
			workflow.Go(putActivityCtx, func(ctx workflow.Context) {
				var putResult *PutActivityResult
				if err := workflow.ExecuteActivity(
					ctx,
					ActivityPutName,
					PutActivityParameters{
						Team:    team,
						Service: service,
					},
				).Get(ctx, &putResult); err != nil {
					logger.Error("error in put", "team", team, "service", service, "error", err)
					putResult = &PutActivityResult{
						Verdict: &models.CheckerVerdict{
							Action:  checkerpb.Action_ACTION_PUT,
							Status:  checkerpb.Status_STATUS_CHECK_FAILED,
							Public:  "internal error",
							Private: fmt.Sprintf("put activity err: %v", err),
						},
					}
				}
				putResultsChan.Send(ctx, putResult)
			})
		}

		getResultsChan := workflow.NewBufferedChannel(ctx, 3)

		for i := 0; i < getCount; i++ {
			workflow.Go(getActivityCtx, func(ctx workflow.Context) {
				var getResult *GetActivityResult
				if err := workflow.ExecuteActivity(
					ctx,
					ActivityGetName,
					GetActivityParameters{
						Team:    team,
						Service: service,
					},
				).Get(ctx, &getResult); err != nil {
					logger.Error("error in get", "team", team, "service", service, "error", err)
					getResult = &GetActivityResult{
						Verdict: &models.CheckerVerdict{
							Action:  checkerpb.Action_ACTION_GET,
							Status:  checkerpb.Status_STATUS_CHECK_FAILED,
							Public:  "internal error",
							Private: fmt.Sprintf("get activity err: %v", err),
						},
					}
				}
				getResultsChan.Send(ctx, getResult)
			})
		}

		for i := 0; i < putCount; i++ {
			var putResult *PutActivityResult
			putResultsChan.Receive(ctx, &putResult)
			putResults = append(putResults, putResult)
		}

		for i := 0; i < getCount; i++ {
			var getResult *GetActivityResult
			getResultsChan.Receive(ctx, &getResult)
			getResults = append(getResults, getResult)
		}
	}

	lastActivityOptions := workflow.ActivityOptions{
		ScheduleToCloseTimeout: time.Second * 5,
	}
	lastCtx := workflow.WithActivityOptions(ctx, lastActivityOptions)
	if err := workflow.ExecuteActivity(
		lastCtx,
		ActivityLastName,
		LastActivityParameters{
			Team:    team,
			Service: service,

			CheckResult: checkResult,
			PutResults:  putResults,
			GetResults:  getResults,
		},
	).Get(ctx, nil); err != nil {
		logger.Error("error in last", "team", team, "service", service, "error", err)
	}

	logger.Info("iteration finished", "team", team, "service", service)
}
