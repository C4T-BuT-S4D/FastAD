package impl

import (
	"context"
	"fmt"

	"go.temporal.io/sdk/activity"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
	"go.temporal.io/sdk/workflow"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/c4t-but-s4d/fastad/internal/checkers"
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

func Run(runCtx, shutdownCtx context.Context, cfg *checkers.Config) error {
	temporalClientOpts := client.Options{
		HostPort: cfg.Temporal.Address,
		Logger: logging.NewTemporalAdapter(
			zap.L().With(zap.String("component", "checkers_worker")),
		),
	}
	if cfg.MetricsAddress != "" {
		handler, err := metrics.TemporalHandler(cfg.MetricsAddress)
		if err != nil {
			return fmt.Errorf("creating metrics handler: %w", err)
		}
		temporalClientOpts.MetricsHandler = handler
	}
	temporalClient, err := client.Dial(temporalClientOpts)
	if err != nil {
		return fmt.Errorf("creating temporal client: %w", err)
	}
	defer temporalClient.Close()

	db := cfg.Postgres.BunDB()

	checkersController := checkers.NewController(db)

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

	checkersWorker := worker.New(temporalClient, "checkers", worker.Options{
		DisableRegistrationAliasing: true,
	})

	// Activities.
	checkersWorker.RegisterActivityWithOptions(
		checkers.NewCheckActivity().ActivityDefinition,
		activity.RegisterOptions{Name: checkers.CheckActivityName},
	)

	checkersWorker.RegisterActivityWithOptions(
		checkers.NewFetchDataActivity(
			teamsClient,
			servicesClient,
			gameStateClient,
		).ActivityDefinition,
		activity.RegisterOptions{Name: checkers.FetchDataActivityName},
	)

	checkersWorker.RegisterActivityWithOptions(
		checkers.NewGetActivity().ActivityDefinition,
		activity.RegisterOptions{Name: checkers.GetActivityName},
	)

	checkersWorker.RegisterActivityWithOptions(
		checkers.NewPickGetFlagActivity(checkersController).ActivityDefinition,
		activity.RegisterOptions{Name: checkers.PickGetFlagActivityName},
	)

	checkersWorker.RegisterActivityWithOptions(
		checkers.NewPrepareRoundActivity(checkersController, gameStateClient).ActivityDefinition,
		activity.RegisterOptions{Name: checkers.PrepareRoundActivityName},
	)

	checkersWorker.RegisterActivityWithOptions(
		checkers.NewPutActivity().ActivityDefinition,
		activity.RegisterOptions{Name: checkers.PutActivityName},
	)

	checkersWorker.RegisterActivityWithOptions(
		checkers.NewSaveRoundDataActivity(checkersController).ActivityDefinition,
		activity.RegisterOptions{Name: checkers.SaveRoundDataActivityName},
	)

	checkersWorker.RegisterActivityWithOptions(
		checkers.NewSaveVerdictActivity(checkersController).ActivityDefinition,
		activity.RegisterOptions{Name: checkers.SaveVerdictActivityName},
	)

	checkersWorker.RegisterActivityWithOptions(
		checkers.NewGetLastExecutionActivity(checkersController).ActivityDefinition,
		activity.RegisterOptions{Name: checkers.GetLastExecutionActivityName},
	)
	// End of activities.

	// Workflows.
	checkersWorker.RegisterWorkflowWithOptions(
		checkers.CheckWorkflowDefinition,
		workflow.RegisterOptions{Name: checkers.CheckWorkflowName},
	)

	checkersWorker.RegisterWorkflowWithOptions(
		checkers.GetWorkflowDefinition,
		workflow.RegisterOptions{Name: checkers.GetWorkflowName},
	)

	checkersWorker.RegisterWorkflowWithOptions(
		checkers.RoundWorkflowDefinition,
		workflow.RegisterOptions{Name: checkers.RoundWorkflowName},
	)
	// End of workflows.

	if cfg.MetricsAddress != "" {
		go metrics.RunServer(runCtx, shutdownCtx, cfg.MetricsAddress)
	}

	go func() {
		<-runCtx.Done()
		checkersWorker.Stop()
	}()
	if err := checkersWorker.Run(nil); err != nil {
		return fmt.Errorf("running worker: %w", err)
	}

	return nil
}
