package checkers

import (
	"context"
	"fmt"

	"github.com/sirupsen/logrus"

	"github.com/c4t-but-s4d/fastad/internal/models"
)

type ActivityFetchDataParameters struct {
	GameState *models.GameState
}

type ActivityFetchDataResult struct {
	Teams    []*models.Team
	Services []*models.Service
}

func (s *ActivityState) ActivityFetchDataDefinition(
	ctx context.Context,
	_ *ActivityFetchDataParameters,
) (*ActivityFetchDataResult, error) {
	teams, err := s.teamsClient.List(ctx)
	if err != nil {
		return nil, fmt.Errorf("getting teams: %w", err)
	}
	logrus.Infof("fetched teams: %v", teams)

	services, err := s.servicesClient.List(ctx)
	if err != nil {
		return nil, fmt.Errorf("getting services: %w", err)
	}
	logrus.Infof("fetched services: %v", services)

	return &ActivityFetchDataResult{
		Teams:    teams,
		Services: services,
	}, nil
}
