package checkers

import (
	"context"

	"github.com/sirupsen/logrus"

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
	logger := logrus.WithFields(logrus.Fields{
		"team":    params.Team.Name,
		"service": params.Service.Name,
		"action":  checkerpb.Action_ACTION_GET,
	})

	logger.Info("starting")
	verdict := RunGetAction(ctx, params)
	logger.Infof("finished: %v", verdict)

	return &GetActivityResult{Verdict: verdict}, nil
}
