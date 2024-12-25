package impl

import (
	"context"
	"fmt"
	"sync"

	"github.com/c4t-but-s4d/fastad/internal/scoreboard"
	"github.com/c4t-but-s4d/fastad/pkg/grpcext"
	"github.com/c4t-but-s4d/fastad/pkg/metrics"
	scoreboardpb "github.com/c4t-but-s4d/fastad/pkg/proto/scoreboard"
)

func Run(runCtx, shutdownCtx context.Context, cfg *scoreboard.Config) error {
	db := cfg.Postgres.BunDB()

	service := scoreboard.NewService(db, cfg)

	if err := service.RestoreState(runCtx); err != nil {
		return fmt.Errorf("restoring state: %w", err)
	}

	wg := sync.WaitGroup{}
	defer wg.Wait()

	wg.Add(1)
	go func() {
		defer wg.Done()
		service.Run(runCtx)
	}()

	grpcServer := grpcext.NewServer(grpcext.WithServerInstallation(cfg.Installation))
	scoreboardpb.RegisterScoreboardServiceServer(grpcServer, service)

	if cfg.MetricsAddress != "" {
		go metrics.RunServer(runCtx, shutdownCtx, cfg.MetricsAddress)
	}

	if err := grpcext.RunServer(runCtx, shutdownCtx, grpcServer, cfg.ListenAddress); err != nil {
		return fmt.Errorf("running server: %w", err)
	}
	return nil
}
