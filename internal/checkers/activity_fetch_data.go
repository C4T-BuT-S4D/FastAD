package checkers

import (
	"github.com/c4t-but-s4d/fastad/internal/models"
)

type ActivityFetchDataParameters struct {
	GameSettings *models.GameSettings
}

type ActivityFetchDataResult struct {
	Teams    []*models.Team
	Services []*models.Service
}

func ActivityFetchDataDefinition(params *ActivityFetchDataParameters) (*ActivityFetchDataResult, error) {
	return &ActivityFetchDataResult{
		Teams: []*models.Team{
			{
				Name:    "team1",
				Address: "addr1",
			},
			{
				Name:    "team2",
				Address: "addr2",
			},
		},
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
