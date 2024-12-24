package main

import (
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"

	"github.com/c4t-but-s4d/fastad/internal/api"
	"github.com/c4t-but-s4d/fastad/pkg/baseconfig"
	"github.com/c4t-but-s4d/fastad/pkg/clients/gamestate"
	"github.com/c4t-but-s4d/fastad/pkg/clients/services"
	"github.com/c4t-but-s4d/fastad/pkg/clients/teams"
	"github.com/c4t-but-s4d/fastad/pkg/grpcext"
	"github.com/c4t-but-s4d/fastad/pkg/httpext"
	"github.com/c4t-but-s4d/fastad/pkg/logging"
	gspb "github.com/c4t-but-s4d/fastad/pkg/proto/data/game_state"
	servicespb "github.com/c4t-but-s4d/fastad/pkg/proto/data/services"
	teamspb "github.com/c4t-but-s4d/fastad/pkg/proto/data/teams"
	receiverpb "github.com/c4t-but-s4d/fastad/pkg/proto/receiver"
	scoreboardpb "github.com/c4t-but-s4d/fastad/pkg/proto/scoreboard"
)

func main() {
	defer logging.Init().Close()

	cfg := baseconfig.MustSetupAll(&api.Config{}, baseconfig.WithEnvPrefix("API"))

	dataServiceConn, err := grpcext.Dial(cfg.DataService.Address, cfg.UserAgent)
	if err != nil {
		zap.L().Fatal("unable to connect to data service", zap.Error(err))
	}

	teamsClient := teams.NewClient(teamspb.NewTeamsServiceClient(dataServiceConn))
	servicesClient := services.NewClient(servicespb.NewServicesServiceClient(dataServiceConn))
	gameStateClient := gamestate.NewClient(gspb.NewGameStateServiceClient(dataServiceConn))

	receiverConn, err := grpcext.Dial(cfg.ReceiverAddress, cfg.UserAgent)
	if err != nil {
		zap.L().Fatal("unable to connect to receiver service", zap.Error(err))
	}
	receiverClient := receiverpb.NewReceiverServiceClient(receiverConn)

	scoreboardConn, err := grpcext.Dial(cfg.ScoreboardAddress, cfg.UserAgent)
	if err != nil {
		zap.L().Fatal("unable to connect to scoreboard service", zap.Error(err))
	}
	scoreboardClient := scoreboardpb.NewScoreboardServiceClient(scoreboardConn)

	apiService := api.NewService(
		teamsClient,
		servicesClient,
		gameStateClient,
		receiverClient,
		scoreboardClient,
	)

	e := echo.New()
	e.HTTPErrorHandler = httpext.ErrorHandler()
	apiService.RegisterRoutes(e)

	if err := e.Start(cfg.ListenAddress); err != nil {
		zap.L().Fatal("unable to start server", zap.Error(err))
	}
}
