package checkers

import (
	"fmt"
	"time"

	"go.temporal.io/sdk/workflow"

	checkerpb "github.com/c4t-but-s4d/fastad/pkg/proto/checker"
)

const RoundWorkflowName = "RoundWorkflow"

type RoundWorkflowParameters struct{}

func RoundWorkflowDefinition(ctx workflow.Context, _ RoundWorkflowParameters) error {
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
		logger.Error("running fetch data activity", "error", err)
		return fmt.Errorf("fetch data: %w", err)
	}

	logger.Info(
		"fetched current data",
		"teams",
		len(fetchDataResult.Teams),
		"services",
		len(fetchDataResult.Services),
		"round",
		fetchDataResult.GameState.RunningRound,
	)

	fetchDataResult.GameState.RunningRound++
	fetchDataResult.GameState.RunningRoundStart = workflow.Now(ctx)

	var prepareStateResult *PrepareRoundActivityResult
	if err := workflow.ExecuteLocalActivity(
		laoCtx,
		ActivityPrepareRoundStateName,
		&PrepareRoundActivityParameters{
			GameState: fetchDataResult.GameState,
			Teams:     fetchDataResult.Teams,
			Services:  fetchDataResult.Services,
		},
	).Get(ctx, &prepareStateResult); err != nil {
		logger.Error("running prepare state activity", "error", err)
		return fmt.Errorf("prepare state: %w", err)
	}

	putsCount := len(fetchDataResult.Teams) * len(fetchDataResult.Services)
	putResultsChan := workflow.NewBufferedChannel(ctx, putsCount)

	wg := workflow.NewWaitGroup(ctx)
	for _, flagInfo := range prepareStateResult.Flags {
		wg.Add(1)
		workflow.Go(ctx, func(ctx workflow.Context) {
			defer wg.Done()

			putActivityCtx := workflow.WithActivityOptions(ctx, workflow.ActivityOptions{
				ScheduleToCloseTimeout: flagInfo.Service.CheckerTimeout(checkerpb.Action_ACTION_PUT) + checkerKillDelay*2,
			})

			var putResult *PutActivityResult
			if err := workflow.ExecuteActivity(
				putActivityCtx,
				ActivityPutName,
				PutActivityParameters{
					FlagInfo: flagInfo,
				},
			).Get(ctx, &putResult); err != nil {
				logger.Error("error in put", "team", flagInfo.Team, "service", flagInfo.Service, "error", err)
				putResult = &PutActivityResult{
					Verdict: &Verdict{
						Action:  checkerpb.Action_ACTION_PUT,
						Status:  checkerpb.Status_STATUS_CHECK_FAILED,
						Public:  "internal error",
						Private: fmt.Sprintf("put activity err: %v", err),
					},
					FlagInfo: flagInfo,
				}
			}

			putResultsChan.Send(ctx, putResult)
		})
	}

	putResults := make([]*PutActivityResult, 0, putsCount)
	for range putsCount {
		var putResult *PutActivityResult
		putResultsChan.Receive(ctx, &putResult)
		putResults = append(putResults, putResult)
	}

	wg.Wait(ctx)

	logger.Info("finished puts, saving results", "put_results", len(putResults))

	var saveRoundDataResult *SaveRoundDataActivityResult
	if err := workflow.ExecuteLocalActivity(
		laoCtx,
		ActivitySaveRoundStateName,
		&SaveRoundDataActivityParameters{
			PutResults: putResults,
		},
	).Get(ctx, &saveRoundDataResult); err != nil {
		logger.Error("running save round data activity", "error", err)
		return fmt.Errorf("save round data: %w", err)
	}

	logger.Info("finished round", "round", fetchDataResult.GameState.RunningRound)

	return nil
}
