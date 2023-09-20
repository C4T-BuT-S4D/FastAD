package checkers

import (
	"context"

	"github.com/c4t-but-s4d/fastad/internal/models"
	"github.com/sirupsen/logrus"
)

type LastActivityParameters struct {
	GameSettings *models.GameSettings
	Team         *models.Team
	Service      *models.Service

	CheckResult *CheckActivityResult
	PutResults  []*PutActivityResult
	GetResults  []*GetActivityResult
}

func LastActivityDefinition(ctx context.Context, params LastActivityParameters) error {
	logrus.Infof("running last %v/%v", params.Team, params.Service)
	logrus.Infof("received check result: %v", params.CheckResult)
	logrus.Infof("received put results: %v", params.PutResults)
	logrus.Infof("received get results: %v", params.GetResults)
	return nil
}
