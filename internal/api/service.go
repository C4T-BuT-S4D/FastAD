package api

import (
	"github.com/labstack/echo/v4"

	"github.com/c4t-but-s4d/fastad/pkg/clients/gamestate"
	"github.com/c4t-but-s4d/fastad/pkg/clients/services"
	"github.com/c4t-but-s4d/fastad/pkg/clients/teams"
	"github.com/c4t-but-s4d/fastad/pkg/httpext"
	receiverpb "github.com/c4t-but-s4d/fastad/pkg/proto/receiver"
	"github.com/c4t-but-s4d/fastad/pkg/proto/scoreboard"
)

type Service struct {
	teamsClient      *teams.Client
	servicesClient   *services.Client
	gameStateClient  *gamestate.Client
	receiverClient   receiverpb.ReceiverServiceClient
	scoreboardClient scoreboard.ScoreboardServiceClient
}

func NewService(
	teamsClient *teams.Client,
	servicesClient *services.Client,
	gameStateClient *gamestate.Client,
	receiverClient receiverpb.ReceiverServiceClient,
	scoreboardClient scoreboard.ScoreboardServiceClient,
) *Service {
	return &Service{
		teamsClient:      teamsClient,
		servicesClient:   servicesClient,
		gameStateClient:  gameStateClient,
		receiverClient:   receiverClient,
		scoreboardClient: scoreboardClient,
	}
}

func (s *Service) RegisterRoutes(e *echo.Echo) {
	apiGroup := e.Group("/api", httpext.RequestIDMiddleware)

	apiGroup.GET("/teams", s.HandleTeamsList())
	apiGroup.GET("/services", s.HandleServicesList())
	apiGroup.GET("/scoreboard", s.HandleGetScoreboard())
	apiGroup.GET("/ctftime", s.HandleGetCTFTimeScoreboard())
	apiGroup.GET("/game_state", s.HandleGetGameState())
}
