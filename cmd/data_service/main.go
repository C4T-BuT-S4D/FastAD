package main

import (
	"context"
	"net"
	"os/signal"
	"syscall"

	"github.com/sirupsen/logrus"

	"github.com/c4t-but-s4d/fastad/internal/config"
	"github.com/c4t-but-s4d/fastad/internal/logging"
	"github.com/c4t-but-s4d/fastad/internal/services/gamestate"
	"github.com/c4t-but-s4d/fastad/internal/services/services"
	"github.com/c4t-but-s4d/fastad/internal/services/teams"
	"github.com/c4t-but-s4d/fastad/internal/version"
	"github.com/c4t-but-s4d/fastad/pkg/grpcext"
	gspb "github.com/c4t-but-s4d/fastad/pkg/proto/data/game_state"
	servicespb "github.com/c4t-but-s4d/fastad/pkg/proto/data/services"
	teamspb "github.com/c4t-but-s4d/fastad/pkg/proto/data/teams"
)

type Config struct {
	ListenAddress string          `mapstructure:"listen_address" default:"127.0.0.1:1337"`
	Postgres      config.Postgres `mapstructure:"postgres"`
}

func main() {
	cfg := config.MustSetupAll(&Config{}, config.WithEnvPrefix("FASTAD_DATA_SERVICE"))

	logging.Init()

	db := cfg.Postgres.BunDB()

	versionController := version.NewController(db)

	teamsController := teams.NewController(db, versionController)
	teamsService := teams.NewService(teamsController)

	servicesController := services.NewController(db, versionController)
	servicesService := services.NewService(servicesController)

	gameStateController := gamestate.NewController(db, versionController)
	gameStateService := gamestate.NewService(gameStateController)

	server := grpcext.NewServer()
	teamspb.RegisterTeamsServiceServer(server, teamsService)
	servicespb.RegisterServicesServiceServer(server, servicesService)
	gspb.RegisterGameStateServiceServer(server, gameStateService)

	runCtx, runCancel := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT)
	defer runCancel()

	logrus.WithField("listen_address", cfg.ListenAddress).Info("starting server")
	lis, err := net.Listen("tcp", cfg.ListenAddress)
	if err != nil {
		logrus.WithError(err).Fatal("error creating listener")
	}

	go func() {
		<-runCtx.Done()
		server.GracefulStop()
	}()

	if err := server.Serve(lis); err != nil {
		logrus.WithError(err).Fatal("error in server")
	}
}
