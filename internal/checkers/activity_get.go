package checkers

import (
	"context"
	"strings"

	"github.com/c4t-but-s4d/fastad/internal/models"
	checkerpb "github.com/c4t-but-s4d/fastad/pkg/proto/checker"
	"github.com/sirupsen/logrus"
)

type GetActivityParameters struct {
	GameSettings *models.GameSettings
	Team         *models.Team
	Service      *models.Service
}

type GetActivityResult struct {
	Verdict *models.CheckerVerdict
}

func GetActivityDefinition(ctx context.Context, params *GetActivityParameters) (*GetActivityResult, error) {
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
	logger.Info("finished: %v", verdict)

	return &GetActivityResult{Verdict: verdict}, nil
}
