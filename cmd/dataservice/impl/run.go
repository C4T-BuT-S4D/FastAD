package impl

import (
	"context"
	"fmt"

	"github.com/c4t-but-s4d/fastad/internal/dataservice"
	"github.com/c4t-but-s4d/fastad/internal/dataservice/gamestate"
	"github.com/c4t-but-s4d/fastad/internal/dataservice/services"
	"github.com/c4t-but-s4d/fastad/internal/dataservice/teams"
	"github.com/c4t-but-s4d/fastad/internal/version"
	"github.com/c4t-but-s4d/fastad/pkg/grpcext"
	"github.com/c4t-but-s4d/fastad/pkg/metrics"
	gspb "github.com/c4t-but-s4d/fastad/pkg/proto/data/game_state"
	servicespb "github.com/c4t-but-s4d/fastad/pkg/proto/data/services"
	teamspb "github.com/c4t-but-s4d/fastad/pkg/proto/data/teams"
)

func Run(runCtx, shutdownCtx context.Context, cfg *dataservice.Config) error {
	db := cfg.Postgres.BunDB()

	versionController := version.NewController(db)

	teamsController := teams.NewController(db, versionController)
	teamsService := teams.NewService(teamsController)

	servicesController := services.NewController(db, versionController)
	servicesService := services.NewService(servicesController)

	gameStateController := gamestate.NewController(db, versionController)
	gameStateService := gamestate.NewService(gameStateController)

	server := grpcext.NewServer(grpcext.WithServerInstallation(cfg.Installation))
	teamspb.RegisterTeamsServiceServer(server, teamsService)
	servicespb.RegisterServicesServiceServer(server, servicesService)
	gspb.RegisterGameStateServiceServer(server, gameStateService)

	if cfg.MetricsAddress != "" {
		go metrics.RunServer(runCtx, shutdownCtx, cfg.MetricsAddress)
	}

	if err := grpcext.RunServer(runCtx, shutdownCtx, server, cfg.ListenAddress); err != nil {
		return fmt.Errorf("running server: %w", err)
	}

	return nil
}
