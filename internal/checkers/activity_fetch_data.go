package checkers

import (
	"context"
	"fmt"

	"github.com/c4t-but-s4d/fastad/internal/models"
	"github.com/sirupsen/logrus"
)

type ActivityFetchDataParameters struct {
	GameSettings *models.GameSettings
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

	return &ActivityFetchDataResult{
		Teams: teams,
		Services: []*models.Service{
			{
				Name: "service1",
			},
			{
				Name: "service2",
			},
		},
	}, nil
}
