package main

import (
	"context"
	"net"
	"os/signal"
	"syscall"

	"go.uber.org/zap"

	"github.com/c4t-but-s4d/fastad/internal/dataservice/gamestate"
	"github.com/c4t-but-s4d/fastad/internal/dataservice/services"
	"github.com/c4t-but-s4d/fastad/internal/dataservice/teams"
	"github.com/c4t-but-s4d/fastad/internal/version"
	"github.com/c4t-but-s4d/fastad/pkg/baseconfig"
	"github.com/c4t-but-s4d/fastad/pkg/config"
	"github.com/c4t-but-s4d/fastad/pkg/grpcext"
	"github.com/c4t-but-s4d/fastad/pkg/logging"
	gspb "github.com/c4t-but-s4d/fastad/pkg/proto/data/game_state"
	servicespb "github.com/c4t-but-s4d/fastad/pkg/proto/data/services"
	teamspb "github.com/c4t-but-s4d/fastad/pkg/proto/data/teams"
)

type Config struct {
	ListenAddress string          `mapstructure:"listen_address" default:"127.0.0.1:1337"`
	Postgres      config.Postgres `mapstructure:"postgres"`
}

func main() {
	defer logging.Init().Close()

	cfg := baseconfig.MustSetupAll(&Config{}, baseconfig.WithEnvPrefix("FASTAD_DATA_SERVICE"))

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

	zap.L().Info("starting server", zap.String("listen_address", cfg.ListenAddress))
	lis, err := net.Listen("tcp", cfg.ListenAddress)
	if err != nil {
		zap.L().Fatal("error creating listener", zap.Error(err))
	}

	go func() {
		<-runCtx.Done()
		server.GracefulStop()
	}()

	if err := server.Serve(lis); err != nil {
		zap.L().Fatal("error in server", zap.Error(err))
	}
}
