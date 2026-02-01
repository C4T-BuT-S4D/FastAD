package checkers

import (
	"context"
	"errors"
	"fmt"

	"github.com/c4t-but-s4d/fastad/internal/models"
	checkerpb "github.com/c4t-but-s4d/fastad/pkg/proto/checker"
)

const GetLastExecutionActivityName = "GetLastExecution"

type GetLastExecutionActivity struct {
	checkersController *Controller
}

func NewGetLastExecutionActivity(checkersController *Controller) *GetLastExecutionActivity {
	return &GetLastExecutionActivity{checkersController: checkersController}
}

type GetLastExecutionActivityParameters struct {
	Action    checkerpb.Action
	TeamID    int
	ServiceID int
}

type GetLastExecutionActivityResult struct {
	Execution *models.CheckerExecution
}

func (a *GetLastExecutionActivity) ActivityDefinition(
	ctx context.Context,
	params *GetLastExecutionActivityParameters,
) (*GetLastExecutionActivityResult, error) {
	lastPutExecution, err := a.checkersController.GetLastExecution(
		ctx,
		params.TeamID,
		params.ServiceID,
		params.Action,
	)
	if errors.Is(err, ErrExecutionNotFound) {
		return &GetLastExecutionActivityResult{Execution: nil}, nil
	}
	if err != nil {
		return nil, fmt.Errorf("getting last execution: %w", err)
	}

	return &GetLastExecutionActivityResult{
		Execution: lastPutExecution,
	}, nil
}
