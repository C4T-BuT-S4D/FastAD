package checkers

import (
	"fmt"
	"time"

	"go.temporal.io/sdk/workflow"

	"github.com/c4t-but-s4d/fastad/internal/models"
	checkerpb "github.com/c4t-but-s4d/fastad/pkg/proto/checker"
)

type WorkflowParameters struct{}

func WorkflowDefinition(ctx workflow.Context, _ WorkflowParameters) error {
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
		&ActivityFetchDataParameters{},
	).Get(ctx, &fetchDataResult); err != nil {
		return fmt.Errorf("running fetch data activity: %w", err)
	}

	wg := workflow.NewWaitGroup(ctx)
	wg.Add(len(fetchDataResult.Teams) * len(fetchDataResult.Services))
	for _, team := range fetchDataResult.Teams {
		for _, service := range fetchDataResult.Services {
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
	getActivityCtx := workflow.WithActivityOptions(ctx, workflow.ActivityOptions{
		ScheduleToCloseTimeout: service.CheckerTimeout(checkerpb.Action_ACTION_GET) + checkerKillDelay*2,
	})

	getCount := service.GetRunCount(checkerpb.Action_ACTION_GET)

	var checkResult *CheckActivityResult
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
			Verdict: &Verdict{
				Action:  checkerpb.Action_ACTION_CHECK,
				Status:  checkerpb.Status_STATUS_CHECK_FAILED,
				Public:  "internal error",
				Private: fmt.Sprintf("check activity err: %v", err),
			},
		}
	}

	if checkResult.Verdict.IsUp() {
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
						Verdict: &Verdict{
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
			GetResults:  getResults,
		},
	).Get(ctx, nil); err != nil {
		logger.Error("error in last", "team", team, "service", service, "error", err)
	}

	logger.Info("iteration finished", "team", team, "service", service)
}
