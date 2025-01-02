package checkers

import (
	"context"
	"fmt"

	"github.com/c4t-but-s4d/fastad/internal/models"
	checkerpb "github.com/c4t-but-s4d/fastad/pkg/proto/checker"
	servicespb "github.com/c4t-but-s4d/fastad/pkg/proto/data/services"
	teamspb "github.com/c4t-but-s4d/fastad/pkg/proto/data/teams"
)

const GetLastExecutionActivityName = "GetLastExecution"

type GetLastExecutionActivity struct {
	checkersController *Controller
}

func NewGetLastExecutionActivity(checkersController *Controller) *GetLastExecutionActivity {
	return &GetLastExecutionActivity{checkersController: checkersController}
}

type GetLastExecutionActivityParameters struct {
	Action  checkerpb.Action
	Team    *teamspb.Team
	Service *servicespb.Service
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
		int(params.Team.Id),
		int(params.Service.Id),
		params.Action,
	)
	if err != nil {
		return nil, fmt.Errorf("getting last execution: %w", err)
	}

	return &GetLastExecutionActivityResult{
		Execution: lastPutExecution,
	}, nil
}
