package checkers

import (
	"context"
	"fmt"

	"github.com/sirupsen/logrus"
)

type CheckActivityParameters struct {
	Team    string
	Service string
}

type CheckActivityResult struct {
	Success bool
}

func CheckActivityDefinition(ctx context.Context, params CheckActivityParameters) (*CheckActivityResult, error) {
	logrus.Infof("running check %s/%s", params.Team, params.Service)
	if params.Team == "team1" && params.Service == "service1" {
		return &CheckActivityResult{Success: true}, nil
	}
	if params.Team == "team2" && params.Service == "service1" {
		return &CheckActivityResult{Success: true}, nil
	}
	return nil, fmt.Errorf("unknown check %s/%s", params.Team, params.Service)
}
