package main

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/c4t-but-s4d/fastad/internal/clients/gamestate"
	"github.com/c4t-but-s4d/fastad/internal/clients/services"
	"github.com/c4t-but-s4d/fastad/internal/clients/teams"
	"github.com/c4t-but-s4d/fastad/internal/config"
	"github.com/c4t-but-s4d/fastad/internal/logging"
	"github.com/c4t-but-s4d/fastad/internal/multiproto"
	"github.com/c4t-but-s4d/fastad/internal/pinger"
	"github.com/c4t-but-s4d/fastad/internal/receiver"
	"github.com/c4t-but-s4d/fastad/pkg/grpcext"
	gspb "github.com/c4t-but-s4d/fastad/pkg/proto/data/game_state"
	servicespb "github.com/c4t-but-s4d/fastad/pkg/proto/data/services"
	teamspb "github.com/c4t-but-s4d/fastad/pkg/proto/data/teams"
	pingerpb "github.com/c4t-but-s4d/fastad/pkg/proto/pinger"
	receiverpb "github.com/c4t-but-s4d/fastad/pkg/proto/receiver"
)

func main() {
	cfg := config.MustSetupAll(&receiver.Config{}, config.WithEnvPrefix("FASTAD_RECEIVER"))

	logging.Init()

	db := cfg.Postgres.BunDB()

	dataServiceConn, err := grpcext.Dial(
		cfg.DataService.Address,
		cfg.UserAgent,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		logrus.WithError(err).Fatal("dialing data service")
	}

	teamsClient := teams.NewClient(teamspb.NewTeamsServiceClient(dataServiceConn))
	servicesClient := services.NewClient(servicespb.NewServicesServiceClient(dataServiceConn))
	gameStateClient := gamestate.NewClient(gspb.NewGameStateServiceClient(dataServiceConn))

	receiverService := receiver.New(db, teamsClient, servicesClient, gameStateClient)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	if err := receiverService.RestoreState(ctx); err != nil {
		logrus.WithError(err).Fatal("restoring receiver state")
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
		logrus.Infof("Running http server on %s", httpServer.Addr)
		if err := httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logrus.WithError(err).Fatal("error running http server")
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c

	shutdownCtx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	if err := httpServer.Shutdown(shutdownCtx); err != nil {
		logrus.WithError(err).Fatal("error shutting down server")
	}
}
