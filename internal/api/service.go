package api

import (
	"context"

	"github.com/centrifugal/centrifuge"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
	"google.golang.org/grpc/metadata"

	"github.com/c4t-but-s4d/fastad/pkg/clients/gamestate"
	"github.com/c4t-but-s4d/fastad/pkg/clients/services"
	"github.com/c4t-but-s4d/fastad/pkg/clients/teams"
	receiverpb "github.com/c4t-but-s4d/fastad/pkg/proto/receiver"
	"github.com/c4t-but-s4d/fastad/pkg/proto/scoreboard"
)

type Service struct {
	centNode *centrifuge.Node

	teamsClient      *teams.Client
	servicesClient   *services.Client
	gameStateClient  *gamestate.Client
	receiverClient   receiverpb.ReceiverServiceClient
	scoreboardClient scoreboard.ScoreboardServiceClient

	boardBuilder *BoardBuilder
}

func NewService(
	centNode *centrifuge.Node,
	teamsClient *teams.Client,
	servicesClient *services.Client,
	gameStateClient *gamestate.Client,
	receiverClient receiverpb.ReceiverServiceClient,
	scoreboardClient scoreboard.ScoreboardServiceClient,
) *Service {
	return &Service{
		centNode: centNode,

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
	apiGroup := e.Group("/api")
	apiGroup.GET("/teams", s.HandleTeamsList())
	apiGroup.GET("/services", s.HandleServicesList())
	apiGroup.GET("/scoreboard", s.HandleGetScoreboard())
	apiGroup.GET("/ctftime", s.HandleGetCTFTimeScoreboard())
	apiGroup.GET("/game_state", s.HandleGetGameState())

	wsHandler := centrifuge.NewWebsocketHandler(s.centNode, centrifuge.WebsocketConfig{})
	e.Any("/centrifuge", echo.WrapHandler(wsHandler), centrifugeAuthMiddleware())
}

func (s *Service) RegisterNode() {
	s.centNode.OnConnect(func(client *centrifuge.Client) {
		transportName := client.Transport().Name()
		transportProto := client.Transport().Protocol()
		zap.L().Debug(
			"client connected",
			zap.String("transport", transportName),
			zap.Any("proto", transportProto),
		)

		client.OnSubscribe(func(event centrifuge.SubscribeEvent, callback centrifuge.SubscribeCallback) {
			zap.L().Debug(
				"client subscribed",
				zap.String("channel", event.Channel),
			)
			// Allow all subscriptions.
			callback(centrifuge.SubscribeReply{}, nil)
		})

		client.OnPublish(func(event centrifuge.PublishEvent, callback centrifuge.PublishCallback) {
			zap.L().Debug(
				"publish event",
				zap.Any("event", event),
			)
			callback(centrifuge.PublishReply{}, nil)
		})

		// Set Disconnect handler to react on client disconnect events.
		client.OnDisconnect(func(e centrifuge.DisconnectEvent) {
			zap.L().Debug(
				"client disconnected",
				zap.Uint32("code", e.Code),
				zap.String("reason", e.Reason),
			)
		})
	})
}

func centrifugeAuthMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			ctx := c.Request().Context()
			cred := &centrifuge.Credentials{
				UserID: "",
			}
			newCtx := centrifuge.SetCredentials(ctx, cred)
			c.SetRequest(c.Request().WithContext(newCtx))
			return next(c)
		}
	}
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
