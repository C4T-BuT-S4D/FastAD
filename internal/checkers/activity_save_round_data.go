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

	logger.Info("saving data for put results", "put_results", len(params.PutResults))

	if err := s.checkersController.SavePutExecutions(ctx, params.PutResults); err != nil {
		return nil, fmt.Errorf("adding checker executions: %w", err)
	}

	return &SaveRoundDataActivityResult{}, nil
}
