package checkers

import (
	"context"

	"github.com/c4t-but-s4d/fastad/internal/models"
	checkerpb "github.com/c4t-but-s4d/fastad/pkg/proto/checker"
	"github.com/sirupsen/logrus"
)

type CheckActivityParameters struct {
	GameSettings *models.GameSettings
	Team         *models.Team
	Service      *models.Service
}

type CheckActivityResult struct {
	Verdict *models.CheckerVerdict
}

func CheckActivityDefinition(ctx context.Context, params *CheckActivityParameters) (*CheckActivityResult, error) {
	logger := logrus.WithFields(logrus.Fields{
		"team":    params.Team.Name,
		"service": params.Service.Name,
		"action":  checkerpb.Action_ACTION_CHECK,
	})

	logger.Info("starting")
	verdict := RunCheckAction(ctx, params)
	logger.Info("finished: %v", verdict)

	return &CheckActivityResult{Verdict: verdict}, nil
}
