package checkers

import (
	"context"
	"fmt"

	"go.temporal.io/sdk/activity"
	"go.temporal.io/sdk/log"

	"github.com/c4t-but-s4d/fastad/internal/models"
	"github.com/c4t-but-s4d/fastad/pkg/clients/gamestate"
	"github.com/c4t-but-s4d/fastad/pkg/clients/services"
	"github.com/c4t-but-s4d/fastad/pkg/clients/teams"
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
	logger := log.With(activity.GetLogger(ctx), "activity", FetchDataActivityName)

	gs, err := a.gameStateClient.Get(ctx)
	if err != nil {
		return nil, fmt.Errorf("getting game state: %w", err)
	}
	logger.Info("fetched game state", "game_state", gs)

	teamsList, err := a.teamsClient.List(ctx)
	if err != nil {
		return nil, fmt.Errorf("getting teams: %w", err)
	}
	logger.Info("fetched teams", "teams", teamsList)

	servicesList, err := a.servicesClient.List(ctx)
	if err != nil {
		return nil, fmt.Errorf("getting services: %w", err)
	}
	logger.Info("fetched services", "services", servicesList)

	return &FetchDataActivityResult{
		GameState: gs,
		Teams:     teamsList,
		Services:  servicesList,
	}, nil
}
