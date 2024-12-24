package impl

import (
	"context"
	"fmt"
	"time"

	"go.temporal.io/sdk/client"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/c4t-but-s4d/fastad/internal/scheduler"
	"github.com/c4t-but-s4d/fastad/pkg/clients/gamestate"
	"github.com/c4t-but-s4d/fastad/pkg/grpcext"
	"github.com/c4t-but-s4d/fastad/pkg/logging"
	gspb "github.com/c4t-but-s4d/fastad/pkg/proto/data/game_state"
)

func Run(runCtx, _ context.Context, cfg *scheduler.Config) error {
	temporalClient, err := client.Dial(client.Options{
		HostPort: cfg.Temporal.Address,
		Logger: logging.NewTemporalAdapter(
			zap.L().With(zap.String("component", "scheduler")),
		),
	})
	if err != nil {
		return fmt.Errorf("creating temporal client: %w", err)
	}
	defer temporalClient.Close()

	db := cfg.Postgres.BunDB()

	dataServiceConn, err := grpcext.Dial(
		cfg.DataService.Address,
		cfg.UserAgent,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return fmt.Errorf("dialing data service: %w", err)
	}

	gameStateClient := gamestate.NewClient(gspb.NewGameStateServiceClient(dataServiceConn))

	t := scheduler.New(time.Second*10, temporalClient, gameStateClient, db)

	if err := t.Run(runCtx); err != nil {
		return fmt.Errorf("running scheduler: %w", err)
	}
	return nil
}
