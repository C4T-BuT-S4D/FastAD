package main

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/c4t-but-s4d/fastad/internal/clients/gamestate"
	gspb "github.com/c4t-but-s4d/fastad/pkg/proto/data/game_state"
	"github.com/creasty/defaults"
	"github.com/mitchellh/mapstructure"
	"github.com/sirupsen/logrus"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
	"go.temporal.io/sdk/activity"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

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

	dataServiceConn, err := grpc.NewClient(
		cfg.DataService.Address,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithUserAgent(cfg.UserAgent),
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

	activityState := checkers.NewActivityState(
		teamsClient,
		servicesClient,
		gameStateClient,
		checkersController,
	)

	checkersWorker := worker.New(temporalClient, "checkers", worker.Options{})
	checkersWorker.RegisterWorkflow(checkers.WorkflowDefinition)

	// Round-related stuff.
	checkersWorker.RegisterWorkflow(checkers.RoundWorkflowDefinition)

	checkersWorker.RegisterActivityWithOptions(
		activityState.PrepareRoundActivityDefinition,
		activity.RegisterOptions{
			Name: checkers.ActivityPrepareRoundStateName,
		},
	)

	checkersWorker.RegisterActivityWithOptions(
		activityState.PutActivityDefinition,
		activity.RegisterOptions{
			Name: checkers.ActivityPutName,
		},
	)

	checkersWorker.RegisterActivityWithOptions(
		activityState.SaveRoundDataActivityDefinition,
		activity.RegisterOptions{
			Name: checkers.ActivitySaveRoundStateName,
		},
	)

	checkersWorker.RegisterActivityWithOptions(
		activityState.ActivityFetchDataDefinition,
		activity.RegisterOptions{
			Name: checkers.ActivityFetchDataName,
		},
	)

	checkersWorker.RegisterActivityWithOptions(
		activityState.CheckActivityDefinition,
		activity.RegisterOptions{
			Name: checkers.ActivityCheckName,
		},
	)

	checkersWorker.RegisterActivityWithOptions(
		activityState.GetActivityDefinition,
		activity.RegisterOptions{
			Name: checkers.ActivityGetName,
		},
	)

	if err := checkersWorker.Run(worker.InterruptCh()); err != nil {
		logrus.WithError(err).Fatalf("error running worker")
	}
}

func setupConfig() (*checkers.Config, error) {
	pflag.BoolP("debug", "v", false, "Enable verbose logging")
	pflag.Parse()

	v := viper.NewWithOptions(viper.ExperimentalBindStruct())

	if err := v.BindPFlags(pflag.CommandLine); err != nil {
		return nil, fmt.Errorf("binding pflags: %w", err)
	}
	v.SetEnvPrefix("FASTAD_CHECKERS")
	v.AutomaticEnv()
	v.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))

	v.SetDefault("user_agent", "checkers_worker")

	cfg := new(checkers.Config)
	defaults.MustSet(cfg)

	if err := v.Unmarshal(
		cfg,
		viper.DecodeHook(
			mapstructure.ComposeDecodeHookFunc(
				mapstructure.TextUnmarshallerHookFunc(),
				mapstructure.StringToTimeDurationHookFunc(),
			),
		),
	); err != nil {
		return nil, fmt.Errorf("unmarshaling config: %w", err)
	}

	return cfg, nil
}
