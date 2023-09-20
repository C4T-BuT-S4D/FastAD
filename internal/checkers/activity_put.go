package checkers

import (
	"context"

	"github.com/sirupsen/logrus"
)

type PutActivityParameters struct {
	Team    string
	Service string
}

type PutActivityResult struct {
	Success bool
}

func PutActivityDefinition(ctx context.Context, params PutActivityParameters) (*PutActivityResult, error) {
	logrus.Infof("running put %s/%s", params.Team, params.Service)
	return nil, nil
}
