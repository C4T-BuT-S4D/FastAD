package impl

import (
	"context"
	"fmt"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

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

	receiverService := receiver.New(db, teamsClient, servicesClient, gameStateClient)

	if err := receiverService.RestoreState(runCtx); err != nil {
		return fmt.Errorf("restoring state: %w", err)
	}

	grpcServer := grpcext.NewServer(grpcext.WithServerInstallation(cfg.Installation))
	receiverpb.RegisterReceiverServiceServer(grpcServer, receiverService)

	if cfg.MetricsAddress != "" {
		go metrics.RunServer(runCtx, shutdownCtx, cfg.MetricsAddress)
	}

	if err := grpcext.RunServer(runCtx, shutdownCtx, grpcServer, cfg.ListenAddress); err != nil {
		return fmt.Errorf("running server: %w", err)
	}
	return nil
}
