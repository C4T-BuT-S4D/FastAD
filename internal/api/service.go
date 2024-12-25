package api

import (
	"context"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
	"google.golang.org/grpc/metadata"

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

	boardBuilder *BoardBuilder
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

		boardBuilder: NewBoardBuilder(
			teamsClient,
			servicesClient,
			receiverClient,
			scoreboardClient,
		),
	}
}

func (s *Service) RegisterRoutes(e *echo.Echo) {
	apiGroup := e.Group("/api", httpext.RequestIDMiddleware())
	apiGroup.GET("/teams", s.HandleTeamsList())
	apiGroup.GET("/services", s.HandleServicesList())
	apiGroup.GET("/scoreboard", s.HandleGetScoreboard())
	apiGroup.GET("/ctftime", s.HandleGetCTFTimeScoreboard())
	apiGroup.GET("/game_state", s.HandleGetGameState())
}

func (s *Service) logger(ctx context.Context) *zap.Logger {
	logger := zap.L().With(zap.String("component", "api"))
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		if requestID, ok := md["request_id"]; ok && len(requestID) > 0 {
			logger = logger.With(zap.String("request_id", requestID[0]))
		}
	}
	return logger
}
