package checkers

import (
	"context"

	"github.com/sirupsen/logrus"
)

type LastActivityParameters struct {
	Team    string
	Service string

	CheckResult *CheckActivityResult
	PutResults  []*PutActivityResult
	GetResults  []*GetActivityResult
}

func LastActivityDefinition(ctx context.Context, params LastActivityParameters) error {
	logrus.Infof("running last %s/%s", params.Team, params.Service)
	logrus.Infof("received check result: %v", params.CheckResult)
	logrus.Infof("received put results: %v", params.PutResults)
	logrus.Infof("received get results: %v", params.GetResults)
	return nil
}
