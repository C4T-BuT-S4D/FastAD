package main

import (
	"context"
	"os/signal"
	"syscall"
	"time"

	"go.temporal.io/sdk/client"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/c4t-but-s4d/fastad/internal/scheduler"
	"github.com/c4t-but-s4d/fastad/pkg/baseconfig"
	"github.com/c4t-but-s4d/fastad/pkg/clients/gamestate"
	"github.com/c4t-but-s4d/fastad/pkg/grpcext"
	"github.com/c4t-but-s4d/fastad/pkg/logging"
	gspb "github.com/c4t-but-s4d/fastad/pkg/proto/data/game_state"
)

func main() {
	defer logging.Init().Close()

	cfg := baseconfig.MustSetupAll(&scheduler.Config{}, baseconfig.WithEnvPrefix("FASTAD_SCHEDULER"))

	temporalClient, err := client.Dial(client.Options{
		HostPort: cfg.Temporal.Address,
		Logger: logging.NewTemporalAdapter(
			zap.L().With(zap.String("component", "scheduler")),
		),
	})
	if err != nil {
		zap.L().Fatal("unable to create temporal client", zap.Error(err))
	}
	defer temporalClient.Close()

	db := cfg.Postgres.BunDB()

	dataServiceConn, err := grpcext.Dial(
		cfg.DataService.Address,
		cfg.UserAgent,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		zap.L().Fatal("unable to connect to data service", zap.Error(err))
	}

	gameStateClient := gamestate.NewClient(gspb.NewGameStateServiceClient(dataServiceConn))

	t := scheduler.New(time.Second*10, temporalClient, gameStateClient, db)

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	if err := t.Run(ctx); err != nil {
		zap.L().Fatal("scheduler run failed", zap.Error(err))
	}
}
