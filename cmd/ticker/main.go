package main

import (
	"context"
	"fmt"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/c4t-but-s4d/fastad/internal/config"
	"github.com/c4t-but-s4d/fastad/internal/logging"
	"github.com/c4t-but-s4d/fastad/internal/ticker"
	"github.com/mitchellh/mapstructure"
	"github.com/sirupsen/logrus"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"go.temporal.io/sdk/client"
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
				"component": "ticker",
			}),
		),
	})
	if err != nil {
		logrus.Fatalf("dialing temporal: %v", err)
	}
	defer temporalClient.Close()

	t := ticker.NewTicker(time.Second*10, temporalClient)

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	t.Run(ctx)
}

func setupConfig() (*ticker.Config, error) {
	pflag.BoolP("debug", "v", false, "Enable verbose logging")
	pflag.Parse()

	if err := viper.BindPFlags(pflag.CommandLine); err != nil {
		return nil, fmt.Errorf("binding pflags: %w", err)
	}

	viper.SetEnvPrefix("FASTAD_TICKER")
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))

	config.SetDefaultPostgresConfig("postgres")
	config.SetDefaultTemporalConfig("temporal")

	cfg := new(ticker.Config)
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
