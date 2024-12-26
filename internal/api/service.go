package api

import (
	"context"

	"github.com/centrifugal/centrifuge"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"

	"github.com/c4t-but-s4d/fastad/internal/centclient"
	"github.com/c4t-but-s4d/fastad/pkg/clients/gamestate"
	"github.com/c4t-but-s4d/fastad/pkg/clients/services"
	"github.com/c4t-but-s4d/fastad/pkg/clients/teams"
	receiverpb "github.com/c4t-but-s4d/fastad/pkg/proto/receiver"
	"github.com/c4t-but-s4d/fastad/pkg/proto/scoreboard"
)

type Service struct {
	config *Config

	centNode         *centrifuge.Node
	teamsClient      *teams.Client
	servicesClient   *services.Client
	gameStateClient  *gamestate.Client
	receiverClient   receiverpb.ReceiverServiceClient
	scoreboardClient scoreboard.ScoreboardServiceClient

	boardBuilder *BoardBuilder
}

func NewService(
	cfg *Config,
	centNode *centrifuge.Node,
	teamsClient *teams.Client,
	servicesClient *services.Client,
	gameStateClient *gamestate.Client,
	receiverClient receiverpb.ReceiverServiceClient,
	scoreboardClient scoreboard.ScoreboardServiceClient,
) *Service {
	return &Service{
		config: cfg,

		centNode:         centNode,
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
	e.Any("/centrifuge", echo.WrapHandler(wsHandler))
}

func (s *Service) RegisterNode() {
	s.centNode.OnConnecting(func(_ context.Context, e centrifuge.ConnectEvent) (centrifuge.ConnectReply, error) {
		// Anonymous user.
		if e.Token == "" {
			return centrifuge.ConnectReply{
				Credentials: &centrifuge.Credentials{
					UserID: "",
				},
			}, nil
		}

		// Intercom user.
		if e.Token != s.config.IntercomToken {
			return centrifuge.ConnectReply{}, centrifuge.DisconnectInvalidToken
		}

		userID, err := centclient.ClientNameFromData(e.Data)
		if err != nil {
			return centrifuge.ConnectReply{}, centrifuge.DisconnectBadRequest
		}

		return centrifuge.ConnectReply{
			Credentials: &centrifuge.Credentials{
				UserID: userID,
			},
		}, nil
	})

	s.centNode.OnConnect(func(client *centrifuge.Client) {
		logger := zap.L().With(
			zap.String("client_id", client.ID()),
			zap.String("user_id", client.UserID()),
		)

		transportName := client.Transport().Name()
		transportProto := client.Transport().Protocol()
		logger.Debug(
			"client connected",
			zap.String("transport", transportName),
			zap.String("proto", string(transportProto)),
		)

		// Allow all subscriptions.
		client.OnSubscribe(func(event centrifuge.SubscribeEvent, callback centrifuge.SubscribeCallback) {
			logger.Debug(
				"client subscribed",
				zap.String("channel", event.Channel),
			)
			callback(centrifuge.SubscribeReply{}, nil)
		})

		// Only allow publishing for authenticated users.
		if client.UserID() != "" {
			client.OnPublish(func(event centrifuge.PublishEvent, callback centrifuge.PublishCallback) {
				logger.Debug(
					"publish event",
					zap.Any("channel", event.Channel),
					zap.ByteString("data", event.Data),
				)
				callback(centrifuge.PublishReply{}, nil)
			})
		}

		client.OnDisconnect(func(e centrifuge.DisconnectEvent) {
			logger.Debug(
				"client disconnected",
				zap.Uint32("code", e.Code),
				zap.String("reason", e.Reason),
			)
		})
	})
}
