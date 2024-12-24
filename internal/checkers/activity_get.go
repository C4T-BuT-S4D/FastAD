package checkers

import (
	"context"

	"go.temporal.io/sdk/activity"
	"go.temporal.io/sdk/log"

	"github.com/c4t-but-s4d/fastad/internal/models"
	checkerpb "github.com/c4t-but-s4d/fastad/pkg/proto/checker"
)

const GetActivityName = "Get"

type GetActivity struct{}

func NewGetActivity() *GetActivity {
	return &GetActivity{}
}

type GetActivityParameters struct {
	GameState *models.GameState
	Team      *models.Team
	Service   *models.Service
	Flag      *models.Flag
}

type GetActivityResult struct {
	Verdict *Verdict
}

func (*GetActivity) ActivityDefinition(ctx context.Context, params *GetActivityParameters) (*GetActivityResult, error) {
	logger := log.With(
		activity.GetLogger(ctx),
		"team", params.Team.Name,
		"service", params.Service.Name,
		"action", checkerpb.Action_ACTION_GET,
		"activity", GetActivityName,
	)

	logger.Info("starting")
	verdict := RunGetAction(ctx, params)
	logger.Info("finished", "verdict", verdict)

	return &GetActivityResult{Verdict: verdict}, nil
}
