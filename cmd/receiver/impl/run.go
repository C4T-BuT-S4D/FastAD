package impl

import (
	"context"
	"fmt"

	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/c4t-but-s4d/fastad/internal/centutil"
	"github.com/c4t-but-s4d/fastad/internal/receiver"
	"github.com/c4t-but-s4d/fastad/pkg/clients/gamestate"
	"github.com/c4t-but-s4d/fastad/pkg/clients/services"
	"github.com/c4t-but-s4d/fastad/pkg/clients/teams"
	"github.com/c4t-but-s4d/fastad/pkg/grpcext"
	"github.com/c4t-but-s4d/fastad/pkg/metrics"
	gspb "github.com/c4t-but-s4d/fastad/pkg/proto/data/game_state"
	servicespb "github.com/c4t-but-s4d/fastad/pkg/proto/data/services"
	teamspb "github.com/c4t-but-s4d/fastad/pkg/proto/data/teams"
	receiverpb "github.com/c4t-but-s4d/fastad/pkg/proto/receiver"
)

func Run(runCtx, shutdownCtx context.Context, cfg *receiver.Config) error {
	db := cfg.Postgres.BunDB()

	producer, err := centutil.NewClientProducer(
		cfg.CentrifugeClient.Address,
		cfg.Channel,
		cfg.Installation,
		cfg.IntercomToken,
	)
	if err != nil {
		return fmt.Errorf("creating centrifuge producer: %w", err)
	}

	dataServiceConn, err := grpcext.Dial(
		cfg.DataService.Address,
		cfg.Installation,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return fmt.Errorf("dialing data service: %w", err)
	}

	teamsClient := teams.NewClient(teamspb.NewTeamsServiceClient(dataServiceConn))
	servicesClient := services.NewClient(servicespb.NewServicesServiceClient(dataServiceConn))
	gameStateClient := gamestate.NewClient(gspb.NewGameStateServiceClient(dataServiceConn))

	receiverService := receiver.New(
		db,
		teamsClient,
		servicesClient,
		gameStateClient,
		producer,
	)

	if err := receiverService.RestoreState(runCtx); err != nil {
		return fmt.Errorf("restoring state: %w", err)
	}

	grpcServer := grpcext.NewServer(grpcext.WithServerInstallation(cfg.Installation))
	receiverpb.RegisterReceiverServiceServer(grpcServer, receiverService)

	g, gctx := errgroup.WithContext(runCtx)

	if cfg.MetricsAddress != "" {
		g.Go(func() error {
			metrics.RunServer(runCtx, shutdownCtx, cfg.MetricsAddress)
			return nil
		})
	}

	g.Go(func() error {
		if err := producer.Run(gctx); err != nil {
			return fmt.Errorf("running client: %w", err)
		}
		return nil
	})

	g.Go(func() error {
		if err := grpcext.RunServer(gctx, shutdownCtx, grpcServer, cfg.ListenAddress); err != nil {
			return fmt.Errorf("running server: %w", err)
		}
		return nil
	})

	if err := g.Wait(); err != nil {
		return fmt.Errorf("waiting: %w", err)
	}

	return nil
}
