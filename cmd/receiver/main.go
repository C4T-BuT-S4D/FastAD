package main

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/c4t-but-s4d/fastad/internal/pinger"
	"github.com/c4t-but-s4d/fastad/internal/receiver"
	"github.com/c4t-but-s4d/fastad/pkg/baseconfig"
	"github.com/c4t-but-s4d/fastad/pkg/clients/gamestate"
	"github.com/c4t-but-s4d/fastad/pkg/clients/services"
	"github.com/c4t-but-s4d/fastad/pkg/clients/teams"
	"github.com/c4t-but-s4d/fastad/pkg/grpcext"
	"github.com/c4t-but-s4d/fastad/pkg/logging"
	"github.com/c4t-but-s4d/fastad/pkg/multiproto"
	gspb "github.com/c4t-but-s4d/fastad/pkg/proto/data/game_state"
	servicespb "github.com/c4t-but-s4d/fastad/pkg/proto/data/services"
	teamspb "github.com/c4t-but-s4d/fastad/pkg/proto/data/teams"
	pingerpb "github.com/c4t-but-s4d/fastad/pkg/proto/pinger"
	receiverpb "github.com/c4t-but-s4d/fastad/pkg/proto/receiver"
)

func main() {
	defer logging.Init().Close()

	cfg := baseconfig.MustSetupAll(&receiver.Config{}, baseconfig.WithEnvPrefix("FASTAD_RECEIVER"))

	db := cfg.Postgres.BunDB()

	dataServiceConn, err := grpcext.Dial(
		cfg.DataService.Address,
		cfg.UserAgent,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		zap.L().Fatal("dialing data service", zap.Error(err))
	}

	teamsClient := teams.NewClient(teamspb.NewTeamsServiceClient(dataServiceConn))
	servicesClient := services.NewClient(servicespb.NewServicesServiceClient(dataServiceConn))
	gameStateClient := gamestate.NewClient(gspb.NewGameStateServiceClient(dataServiceConn))

	receiverService := receiver.New(db, teamsClient, servicesClient, gameStateClient)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	if err := receiverService.RestoreState(ctx); err != nil {
		zap.L().Fatal("restoring receiver state", zap.Error(err))
	}

	grpcServer := grpcext.NewServer()
	receiverpb.RegisterReceiverServiceServer(grpcServer, receiverService)
	pingerpb.RegisterPingerServiceServer(grpcServer, pinger.New())

	httpServer := &http.Server{
		Addr:              "0.0.0.0:8002",
		Handler:           multiproto.NewHandler(grpcServer),
		ReadHeaderTimeout: time.Second * 10,
	}

	go func() {
		zap.L().Info("Running http server", zap.String("listen_address", httpServer.Addr))
		if err := httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			zap.L().Fatal("error running http server", zap.Error(err))
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c

	shutdownCtx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	if err := httpServer.Shutdown(shutdownCtx); err != nil {
		zap.L().Fatal("error shutting down server", zap.Error(err))
	}
}
