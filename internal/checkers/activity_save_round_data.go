package checkers

import (
	"context"
	"fmt"

	"github.com/sirupsen/logrus"

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

func (s *SaveRoundDataActivity) SaveRoundDataActivityDefinition(ctx context.Context, params *SaveRoundDataActivityParameters) (*SaveRoundDataActivityResult, error) {
	logger := logrus.WithFields(logrus.Fields{
		"action": "SaveRoundData",
	})

	logger.Info("starting")
	err := s.saveRoundData(ctx, params, logger)
	if err != nil {
		return nil, fmt.Errorf("saving round data: %w", err)
	}
	logger.Infof("finished")

	return &SaveRoundDataActivityResult{}, nil
}

func (s *SaveRoundDataActivity) saveRoundData(
	ctx context.Context,
	params *SaveRoundDataActivityParameters,
	logger *logrus.Entry,
) error {
	logger.Infof("saving data for %d put results", len(params.PutResults))

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

	logger.Infof("saving %d executions", len(executions))
	if err := s.checkersController.AddCheckerExecutions(ctx, executions); err != nil {
		return fmt.Errorf("adding checker executions: %w", err)
	}

	return nil
}
