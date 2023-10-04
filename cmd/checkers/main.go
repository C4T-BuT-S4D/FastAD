package main

import (
	"context"
	"fmt"
	"os/signal"
	"strings"
	"syscall"

	"github.com/c4t-but-s4d/fastad/internal/checkers"
	"github.com/c4t-but-s4d/fastad/internal/clients/services"
	"github.com/c4t-but-s4d/fastad/internal/clients/teams"
	"github.com/c4t-but-s4d/fastad/internal/config"
	"github.com/c4t-but-s4d/fastad/internal/logging"
	servicespb "github.com/c4t-but-s4d/fastad/pkg/proto/data/services"
	teamspb "github.com/c4t-but-s4d/fastad/pkg/proto/data/teams"
	"github.com/mitchellh/mapstructure"
	"github.com/sirupsen/logrus"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"go.temporal.io/sdk/activity"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	cfg, err := setupConfig()
	if err != nil {
		logrus.Fatalf("error setting up config: %v", err)
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
		logrus.Fatalf("unable to create temporal client: %v", err)
	}
	defer temporalClient.Close()

	runCtx, runCancel := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT)
	defer runCancel()

	dataServiceConn, err := grpc.DialContext(
		runCtx,
		cfg.DataService.Address,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithUserAgent(cfg.UserAgent),
	)
	if err != nil {
		logrus.Fatalf("unable to connect to data service: %v", err)
	}

	teamsClient := teams.NewClient(teamspb.NewTeamsServiceClient(dataServiceConn))

	servicesClient := services.NewClient(servicespb.NewServicesServiceClient(dataServiceConn))

	activityState := checkers.NewActivityState(teamsClient, servicesClient)

	checkersWorker := worker.New(temporalClient, "checkers", worker.Options{})
	checkersWorker.RegisterWorkflow(checkers.WorkflowDefinition)

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
		activityState.PutActivityDefinition,
		activity.RegisterOptions{
			Name: checkers.ActivityPutName,
		},
	)

	checkersWorker.RegisterActivityWithOptions(
		activityState.GetActivityDefinition,
		activity.RegisterOptions{
			Name: checkers.ActivityGetName,
		},
	)

	checkersWorker.RegisterActivityWithOptions(
		activityState.CheckActivityDefinition,
		activity.RegisterOptions{
			Name: checkers.ActivityLastName,
		},
	)

	if err := checkersWorker.Run(worker.InterruptCh()); err != nil {
		logrus.Fatalf("Unable to start workers: %v", err)
	}
}

func setupConfig() (*checkers.Config, error) {
	pflag.BoolP("debug", "v", false, "Enable verbose logging")
	pflag.Parse()
	if err := viper.BindPFlags(pflag.CommandLine); err != nil {
		return nil, fmt.Errorf("binding pflags: %w", err)
	}
	viper.SetEnvPrefix("FASTAD_CHECKERS")
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))

	config.SetDefaultTemporalConfig("temporal")
	config.SetDefaultDataServiceConfig("data_service")

	viper.SetDefault("user_agent", "checkers_worker")

	cfg := new(checkers.Config)
	if err := viper.Unmarshal(
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
