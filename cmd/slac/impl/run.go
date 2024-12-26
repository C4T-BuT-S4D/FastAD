package impl

import (
	"context"
	"fmt"

	"golang.org/x/sync/errgroup"

	"github.com/c4t-but-s4d/fastad/internal/centclient"
	"github.com/c4t-but-s4d/fastad/internal/slac"
	"github.com/c4t-but-s4d/fastad/pkg/grpcext"
	"github.com/c4t-but-s4d/fastad/pkg/metrics"
	slacpb "github.com/c4t-but-s4d/fastad/pkg/proto/slac"
)

func Run(runCtx, shutdownCtx context.Context, cfg *slac.Config) error {
	db := cfg.Postgres.BunDB()

	producer, err := centclient.NewProducer(
		cfg.CentrifugeClient.Address,
		cfg.Channel,
		cfg.Installation,
		cfg.IntercomToken,
	)
	if err != nil {
		return fmt.Errorf("creating centrifuge producer: %w", err)
	}

	service := slac.NewService(db, cfg, producer)

	if err := service.RestoreState(runCtx); err != nil {
		return fmt.Errorf("restoring state: %w", err)
	}

	g, gctx := errgroup.WithContext(runCtx)

	grpcServer := grpcext.NewServer(grpcext.WithServerInstallation(cfg.Installation))
	slacpb.RegisterSlacServiceServer(grpcServer, service)

	if cfg.MetricsAddress != "" {
		g.Go(func() error {
			metrics.RunServer(gctx, shutdownCtx, cfg.MetricsAddress)
			return nil
		})
	}

	g.Go(func() error {
		service.Run(gctx)
		return nil
	})

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
