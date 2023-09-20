package checkers

import (
	"fmt"
	"time"

	"github.com/c4t-but-s4d/fastad/internal/models"
	checkerpb "github.com/c4t-but-s4d/fastad/pkg/proto/checker"
	"go.temporal.io/sdk/workflow"
)

type WorkflowParameters struct {
	GameSettings *models.GameSettings
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
		ActivityFetchDataDefinition,
		&ActivityFetchDataParameters{
			GameSettings: params.GameSettings,
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

	commonActivityOptions := workflow.ActivityOptions{
		ScheduleToCloseTimeout: service.CheckerTimeout() + checkerKillDelay*2,
	}
	activityCtx := workflow.WithActivityOptions(ctx, commonActivityOptions)

	var checkResult *CheckActivityResult
	putResults := make([]*PutActivityResult, 0, service.Puts)
	getResults := make([]*GetActivityResult, 0, service.Gets)

	if err := workflow.ExecuteActivity(activityCtx, CheckActivityDefinition, CheckActivityParameters{
		Team:    team,
		Service: service,
	}).Get(ctx, &checkResult); err != nil {
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
		putResultsChan := workflow.NewBufferedChannel(ctx, service.Puts)

		for i := 0; i < service.Puts; i++ {
			workflow.Go(activityCtx, func(ctx workflow.Context) {
				var putResult *PutActivityResult
				if err := workflow.ExecuteActivity(
					ctx,
					PutActivityDefinition,
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

		for i := 0; i < service.Gets; i++ {
			workflow.Go(activityCtx, func(ctx workflow.Context) {
				var getResult *GetActivityResult
				if err := workflow.ExecuteActivity(
					ctx,
					GetActivityDefinition,
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

		for i := 0; i < service.Puts; i++ {
			var putResult *PutActivityResult
			putResultsChan.Receive(ctx, &putResult)
			putResults = append(putResults, putResult)
		}

		for i := 0; i < service.Gets; i++ {
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
		LastActivityDefinition,
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
