package impl

import (
	"context"
	"fmt"

	"go.temporal.io/sdk/client"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"

	"github.com/c4t-but-s4d/fastad/internal/scheduler"
	"github.com/c4t-but-s4d/fastad/pkg/clients/gamestate"
	"github.com/c4t-but-s4d/fastad/pkg/clients/services"
	"github.com/c4t-but-s4d/fastad/pkg/clients/teams"
	"github.com/c4t-but-s4d/fastad/pkg/grpcext"
	"github.com/c4t-but-s4d/fastad/pkg/logging"
	"github.com/c4t-but-s4d/fastad/pkg/metrics"
	gspb "github.com/c4t-but-s4d/fastad/pkg/proto/data/game_state"
	servicespb "github.com/c4t-but-s4d/fastad/pkg/proto/data/services"
	teamspb "github.com/c4t-but-s4d/fastad/pkg/proto/data/teams"
)

func Run(runCtx, shutdownCtx context.Context, cfg *scheduler.Config) error {
	logger := zap.L().Named(cfg.Installation)

	temporalClient, err := client.Dial(client.Options{
		HostPort: cfg.Temporal.Address,
		Logger: logging.NewTemporalAdapter(
			logger.Named("temporal_client"),
		),
	})
	if err != nil {
		return fmt.Errorf("creating temporal client: %w", err)
	}
	defer temporalClient.Close()

	db := cfg.Postgres.BunDB()

	dataServiceConn, err := grpcext.Dial(
		cfg.DataService.Address,
		cfg.Installation,
		grpcext.AuthDialOptions(cfg.IntercomToken)...,
	)
	if err != nil {
		return fmt.Errorf("dialing data service: %w", err)
	}

	gameStateClient := gamestate.NewClient(gspb.NewGameStateServiceClient(dataServiceConn), cfg.Installation)
	teamsClient := teams.NewClient(teamspb.NewTeamsServiceClient(dataServiceConn), cfg.Installation)
	servicesClient := services.NewClient(servicespb.NewServicesServiceClient(dataServiceConn), cfg.Installation)

	t := scheduler.NewRoundScheduler(
		temporalClient,
		gameStateClient,
		db,
		logger,
	)

	cm := scheduler.NewCheckManager(
		db,
		gameStateClient,
		teamsClient,
		servicesClient,
		temporalClient,
		logger,
	)

	g, gctx := errgroup.WithContext(runCtx)
	g.Go(func() error {
		if err := t.Run(gctx); err != nil {
			return fmt.Errorf("running round scheduler: %w", err)
		}
		return nil
	})

	g.Go(func() error {
		if err := cm.Run(gctx); err != nil {
			return fmt.Errorf("running check manager: %w", err)
		}
		return nil
	})

	if cfg.MetricsAddress != "" {
		g.Go(func() error {
			metrics.RunServer(gctx, shutdownCtx, cfg.MetricsAddress)
			return nil
		})
	}

	if err := g.Wait(); err != nil {
		return fmt.Errorf("running errgroup: %w", err)
	}

	return nil
}
