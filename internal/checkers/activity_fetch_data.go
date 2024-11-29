package checkers

import (
	"context"
	"fmt"

	"github.com/sirupsen/logrus"

	"github.com/c4t-but-s4d/fastad/internal/clients/gamestate"
	"github.com/c4t-but-s4d/fastad/internal/clients/services"
	"github.com/c4t-but-s4d/fastad/internal/clients/teams"
	"github.com/c4t-but-s4d/fastad/internal/models"
)

const FetchDataActivityName = "FetchData"

type FetchDataActivity struct {
	teamsClient     *teams.Client
	servicesClient  *services.Client
	gameStateClient *gamestate.Client
}

func NewFetchDataActivity(
	teamsClient *teams.Client,
	servicesClient *services.Client,
	gameStateClient *gamestate.Client,
) *FetchDataActivity {
	return &FetchDataActivity{
		teamsClient:     teamsClient,
		servicesClient:  servicesClient,
		gameStateClient: gameStateClient,
	}
}

type FetchDataActivityParameters struct{}

type FetchDataActivityResult struct {
	GameState *models.GameState
	Teams     []*models.Team
	Services  []*models.Service
}

func (a *FetchDataActivity) ActivityDefinition(
	ctx context.Context,
	_ *FetchDataActivityParameters,
) (*FetchDataActivityResult, error) {
	gs, err := a.gameStateClient.Get(ctx)
	if err != nil {
		return nil, fmt.Errorf("getting game state: %w", err)
	}
	logrus.Infof("fetched game state: %v", gs)

	teamsList, err := a.teamsClient.List(ctx)
	if err != nil {
		return nil, fmt.Errorf("getting teams: %w", err)
	}
	logrus.Infof("fetched teams: %v", teamsList)

	servicesList, err := a.servicesClient.List(ctx)
	if err != nil {
		return nil, fmt.Errorf("getting services: %w", err)
	}
	logrus.Infof("fetched services: %v", servicesList)

	return &FetchDataActivityResult{
		GameState: gs,
		Teams:     teamsList,
		Services:  servicesList,
	}, nil
}
