package checkers

import (
	"context"
	"strings"

	"github.com/sirupsen/logrus"

	"github.com/c4t-but-s4d/fastad/internal/models"
	checkerpb "github.com/c4t-but-s4d/fastad/pkg/proto/checker"
)

type GetActivityParameters struct {
	GameState *models.GameState
	Team      *models.Team
	Service   *models.Service
}

type GetActivityResult struct {
	Verdict *models.CheckerVerdict
}

func (s *ActivityState) GetActivityDefinition(ctx context.Context, params *GetActivityParameters) (*GetActivityResult, error) {
	logger := logrus.WithFields(logrus.Fields{
		"team":    params.Team.Name,
		"service": params.Service.Name,
		"action":  checkerpb.Action_ACTION_GET,
	})

	flag := &models.Flag{
		Flag:    strings.Repeat("A", 31) + "=",
		Private: "some-flag-id",
	}

	logger.Info("starting")
	verdict := RunGetAction(ctx, params, flag)
	logger.Infof("finished: %v", verdict)

	return &GetActivityResult{Verdict: verdict}, nil
}
