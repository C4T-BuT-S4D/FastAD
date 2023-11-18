package checkers

import (
	"context"

	"github.com/sirupsen/logrus"

	"github.com/c4t-but-s4d/fastad/internal/models"
	checkerpb "github.com/c4t-but-s4d/fastad/pkg/proto/checker"
)

type CheckActivityParameters struct {
	GameState *models.GameState
	Team      *models.Team
	Service   *models.Service
}

type CheckActivityResult struct {
	Verdict *models.CheckerVerdict
}

func (s *ActivityState) CheckActivityDefinition(ctx context.Context, params *CheckActivityParameters) (*CheckActivityResult, error) {
	logger := logrus.WithFields(logrus.Fields{
		"team":    params.Team.Name,
		"service": params.Service.Name,
		"action":  checkerpb.Action_ACTION_CHECK,
	})

	logger.Info("starting")
	verdict := RunCheckAction(ctx, params)
	logger.Infof("finished: %v", verdict)

	return &CheckActivityResult{Verdict: verdict}, nil
}
