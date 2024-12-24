package checkers

import (
	"context"
	"fmt"

	"go.temporal.io/sdk/activity"
	"go.temporal.io/sdk/log"

	"github.com/c4t-but-s4d/fastad/internal/models"
	checkerpb "github.com/c4t-but-s4d/fastad/pkg/proto/checker"
)

const SaveRoundDataActivityName = "SaveRoundData"

type SaveRoundDataActivity struct {
	checkersController *Controller
}

func NewSaveRoundDataActivity(checkersController *Controller) *SaveRoundDataActivity {
	return &SaveRoundDataActivity{
		checkersController: checkersController,
	}
}

type SaveRoundDataActivityParameters struct {
	PutResults []*PutActivityResult
}

type SaveRoundDataActivityResult struct{}

func (s *SaveRoundDataActivity) ActivityDefinition(ctx context.Context, params *SaveRoundDataActivityParameters) (*SaveRoundDataActivityResult, error) {
	logger := log.With(
		activity.GetLogger(ctx),
		"activity", SaveRoundDataActivityName,
	)

	logger.Info("starting")
	err := s.saveRoundData(ctx, params, logger)
	if err != nil {
		return nil, fmt.Errorf("saving round data: %w", err)
	}
	logger.Info("finished")

	return &SaveRoundDataActivityResult{}, nil
}

func (s *SaveRoundDataActivity) saveRoundData(
	ctx context.Context,
	params *SaveRoundDataActivityParameters,
	logger log.Logger,
) error {
	logger.Info("saving data for put results", "put_results", len(params.PutResults))

	executions := make([]*models.CheckerExecution, 0, len(params.PutResults))
	for _, putResult := range params.PutResults {
		execution := &models.CheckerExecution{
			ExecutionID: fmt.Sprintf("put-flag-%d", putResult.FlagInfo.Flag.ID),
			TeamID:      putResult.FlagInfo.Team.ID,
			ServiceID:   putResult.FlagInfo.Service.ID,
			Action:      checkerpb.Action_ACTION_PUT,
			Status:      putResult.Verdict.Status,
			Public:      putResult.Verdict.Public,
			Private:     putResult.Verdict.Private,
			Command:     putResult.Verdict.Command,
		}
		executions = append(executions, execution)
	}

	if len(executions) == 0 {
		logger.Warn("no executions to save")
		return nil
	}

	logger.Info("saving executions", "executions", len(executions))
	if err := s.checkersController.AddCheckerExecutions(ctx, executions); err != nil {
		return fmt.Errorf("adding checker executions: %w", err)
	}

	return nil
}
