package checkers

import (
	"context"
	"fmt"

	"go.temporal.io/sdk/activity"
	"go.temporal.io/sdk/log"
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

	if err := s.checkersController.SavePutExecutions(ctx, params.PutResults); err != nil {
		return fmt.Errorf("adding checker executions: %w", err)
	}

	return nil
}
