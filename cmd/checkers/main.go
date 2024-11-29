package main

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/spf13/pflag"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
	"go.temporal.io/sdk/activity"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
	"go.temporal.io/sdk/workflow"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/c4t-but-s4d/fastad/internal/clients/gamestate"
	"github.com/c4t-but-s4d/fastad/internal/config"
	"github.com/c4t-but-s4d/fastad/pkg/grpctools"
	gspb "github.com/c4t-but-s4d/fastad/pkg/proto/data/game_state"
	"github.com/c4t-but-s4d/fastad/pkg/util"

	"github.com/c4t-but-s4d/fastad/internal/checkers"
	"github.com/c4t-but-s4d/fastad/internal/clients/services"
	"github.com/c4t-but-s4d/fastad/internal/clients/teams"
	"github.com/c4t-but-s4d/fastad/internal/logging"
	servicespb "github.com/c4t-but-s4d/fastad/pkg/proto/data/services"
	teamspb "github.com/c4t-but-s4d/fastad/pkg/proto/data/teams"
)

func main() {
	cfg, err := setupConfig()
	if err != nil {
		logrus.WithError(err).Fatal("error setting up config")
	}

	logging.Init()

	temporalClient, err := client.Dial(client.Options{
		HostPort: cfg.Temporal.Address,
		Logger: logging.NewTemporalAdapter(
			logrus.WithFields(logrus.Fields{
				"component": "checkers_worker",
			}),
		),
	})
	if err != nil {
		logrus.WithError(err).Fatal("unable to create temporal client")
	}
	defer temporalClient.Close()

	pgConn := pgdriver.NewConnector(
		pgdriver.WithAddr(fmt.Sprintf("%s:%d", cfg.Postgres.Host, cfg.Postgres.Port)),
		pgdriver.WithDatabase(cfg.Postgres.Database),
		pgdriver.WithUser(cfg.Postgres.User),
		pgdriver.WithPassword(cfg.Postgres.Password),
		pgdriver.WithInsecure(!cfg.Postgres.EnableSSL),
	)

	sqlDB := sql.OpenDB(pgConn)
	sqlDB.SetMaxIdleConns(cfg.Postgres.MaxIdleConns)
	sqlDB.SetMaxOpenConns(cfg.Postgres.MaxOpenConns)
	sqlDB.SetConnMaxIdleTime(cfg.Postgres.ConnMaxIdleTime)
	sqlDB.SetConnMaxLifetime(cfg.Postgres.ConnMaxLifetime)

	db := bun.NewDB(sqlDB, pgdialect.New())
	logging.AddBunQueryHook(db)

	dataServiceConn, err := grpctools.Dial(
		cfg.DataService.Address,
		cfg.UserAgent,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		logrus.WithError(err).Fatal("unable to connect to data service")
	}

	initCtx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	checkersController := checkers.NewController(db)
	if err := checkersController.MigrateDB(initCtx); err != nil {
		logrus.WithError(err).Fatal("error migrating db")
	}

	teamsClient := teams.NewClient(teamspb.NewTeamsServiceClient(dataServiceConn))
	servicesClient := services.NewClient(servicespb.NewServicesServiceClient(dataServiceConn))
	gameStateClient := gamestate.NewClient(gspb.NewGameStateServiceClient(dataServiceConn))

	checkersWorker := worker.New(temporalClient, "checkers", worker.Options{})

	// Activities.
	checkersWorker.RegisterActivityWithOptions(
		checkers.NewCheckActivity(),
		activity.RegisterOptions{Name: checkers.CheckActivityName},
	)

	checkersWorker.RegisterActivityWithOptions(
		checkers.NewFetchDataActivity(
			teamsClient,
			servicesClient,
			gameStateClient,
		),
		activity.RegisterOptions{Name: checkers.FetchDataActivityName},
	)

	checkersWorker.RegisterActivityWithOptions(
		checkers.NewGetActivity(),
		activity.RegisterOptions{Name: checkers.GetActivityName},
	)

	checkersWorker.RegisterActivityWithOptions(
		checkers.NewPickGetFlagActivity(checkersController),
		activity.RegisterOptions{Name: checkers.PickGetFlagActivityName},
	)

	checkersWorker.RegisterActivityWithOptions(
		checkers.NewPrepareRoundActivity(checkersController, gameStateClient),
		activity.RegisterOptions{Name: checkers.PrepareRoundActivityName},
	)

	checkersWorker.RegisterActivityWithOptions(
		checkers.NewPutActivity(),
		activity.RegisterOptions{Name: checkers.PutActivityName},
	)

	checkersWorker.RegisterActivityWithOptions(
		checkers.NewSaveRoundDataActivity(checkersController),
		activity.RegisterOptions{Name: checkers.SaveRoundDataActivityName},
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

	if err := checkersWorker.Run(worker.InterruptCh()); err != nil {
		logrus.WithError(err).Fatalf("error running worker")
	}
}

func setupConfig() (*checkers.Config, error) {
	pflag.BoolP("debug", "v", false, "Enable verbose logging")
	pflag.Parse()

	cfg := util.Must[*checkers.Config]("setup config")(
		config.SetupAll[*checkers.Config]("FASTAD_CHECKERS"),
	)

	return cfg, nil
}
