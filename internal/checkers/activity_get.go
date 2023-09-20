package checkers

import (
	"context"

	"github.com/sirupsen/logrus"
)

type GetActivityParameters struct {
	Team    string
	Service string
}

type GetActivityResult struct {
	Success bool
}

func GetActivityDefinition(ctx context.Context, params GetActivityParameters) (*GetActivityResult, error) {
	logrus.Infof("running get %s/%s", params.Team, params.Service)
	return nil, nil
}
