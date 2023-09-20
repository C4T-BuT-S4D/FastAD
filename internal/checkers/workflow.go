package checkers

import (
	"time"

	"go.temporal.io/sdk/workflow"
)

type WorkflowParameters struct {
	Teams    []string
	Services []string
}

func WorkflowDefinition(ctx workflow.Context, params WorkflowParameters) error {
	logger := workflow.GetLogger(ctx)
	logger.Info("starting workflow")

	wg := workflow.NewWaitGroup(ctx)
	wg.Add(len(params.Teams) * len(params.Services))
	for _, team := range params.Teams {
		for _, service := range params.Services {
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

func runCheckers(ctx workflow.Context, team string, service string) {
	logger := workflow.GetLogger(ctx)

	checkActivityOptions := workflow.ActivityOptions{
		ScheduleToCloseTimeout: time.Second,
	}
	checkCtx := workflow.WithActivityOptions(ctx, checkActivityOptions)
	var checkResult *CheckActivityResult
	if err := workflow.ExecuteActivity(checkCtx, CheckActivityDefinition, CheckActivityParameters{
		Team:    team,
		Service: service,
	}).Get(ctx, &checkResult); err != nil {
		logger.Error("error in check", "team", team, "service", service, "error", err)
		return
	}

	if !checkResult.Success {
		logger.Info("check failed", "team", team, "service", service)
		return
	}

	putResultsChan := workflow.NewBufferedChannel(ctx, 3)

	putActivityOptions := workflow.ActivityOptions{
		ScheduleToCloseTimeout: time.Second,
	}
	putCtx := workflow.WithActivityOptions(ctx, putActivityOptions)
	for i := 0; i < 3; i++ {
		workflow.Go(putCtx, func(ctx workflow.Context) {
			var putResult *PutActivityResult
			if err := workflow.ExecuteActivity(
				putCtx,
				PutActivityDefinition,
				PutActivityParameters{
					Team:    team,
					Service: service,
				},
			); err != nil {
				logger.Error("error in put", "team", team, "service", service, "error", err)
				putResult = &PutActivityResult{Success: false}
			}
			putResultsChan.Send(ctx, putResult)
		})
	}

	getResultsChan := workflow.NewBufferedChannel(ctx, 3)

	getActivityOptions := workflow.ActivityOptions{
		ScheduleToCloseTimeout: time.Second,
	}
	getCtx := workflow.WithActivityOptions(ctx, getActivityOptions)
	for i := 0; i < 2; i++ {
		workflow.Go(getCtx, func(ctx workflow.Context) {
			var getResult *GetActivityResult
			if err := workflow.ExecuteActivity(
				getCtx,
				GetActivityDefinition,
				GetActivityParameters{
					Team:    team,
					Service: service,
				},
			); err != nil {
				logger.Error("error in get", "team", team, "service", service, "error", err)
				getResult = &GetActivityResult{Success: false}
			}
			getResultsChan.Send(ctx, getResult)
		})
	}

	putResults := make([]*PutActivityResult, 0, 3)
	for i := 0; i < 3; i++ {
		var putResult *PutActivityResult
		putResultsChan.Receive(ctx, &putResult)
		putResults = append(putResults, putResult)
	}

	getResults := make([]*GetActivityResult, 0, 2)
	for i := 0; i < 2; i++ {
		var getResult *GetActivityResult
		getResultsChan.Receive(ctx, &getResult)
		getResults = append(getResults, getResult)
	}

	lastActivityOptions := workflow.ActivityOptions{
		ScheduleToCloseTimeout: time.Second,
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
	); err != nil {
		logger.Error("error in last", "team", team, "service", service, "error", err)
	}

	logger.Info("iteration finished", "team", team, "service", service)
}
