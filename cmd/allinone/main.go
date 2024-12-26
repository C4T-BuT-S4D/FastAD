package main

import (
	"fmt"

	"github.com/google/uuid"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"

	apiImpl "github.com/c4t-but-s4d/fastad/cmd/api/impl"
	checkersImpl "github.com/c4t-but-s4d/fastad/cmd/checkers/impl"
	dataserviceImpl "github.com/c4t-but-s4d/fastad/cmd/dataservice/impl"
	receiverImpl "github.com/c4t-but-s4d/fastad/cmd/receiver/impl"
	schedulerImpl "github.com/c4t-but-s4d/fastad/cmd/scheduler/impl"
	scoreboardImpl "github.com/c4t-but-s4d/fastad/cmd/scoreboard/impl"
	"github.com/c4t-but-s4d/fastad/internal/api"
	"github.com/c4t-but-s4d/fastad/internal/checkers"
	"github.com/c4t-but-s4d/fastad/internal/dataservice"
	"github.com/c4t-but-s4d/fastad/internal/receiver"
	"github.com/c4t-but-s4d/fastad/internal/scheduler"
	"github.com/c4t-but-s4d/fastad/internal/scoreboard"
	"github.com/c4t-but-s4d/fastad/pkg/apiwait"
	"github.com/c4t-but-s4d/fastad/pkg/baseconfig"
	"github.com/c4t-but-s4d/fastad/pkg/config"
	"github.com/c4t-but-s4d/fastad/pkg/logging"
	"github.com/c4t-but-s4d/fastad/pkg/metrics"
	"github.com/c4t-but-s4d/fastad/pkg/stop"
)

type Config struct {
	UserAgent string `mapstructure:"user_agent" default:"allinone"`

	MetricsAddress string `mapstructure:"metrics_address" default:":3000"`
	IntercomToken  string `mapstructure:"intercom_token"`

	Postgres    config.Postgres    `mapstructure:"postgres"`
	Temporal    config.Temporal    `mapstructure:"temporal"`
	DataService dataservice.Config `mapstructure:"data_service"`
	Scheduler   scheduler.Config   `mapstructure:"scheduler"`
	Scoreboard  scoreboard.Config  `mapstructure:"scoreboard"`
	API         api.Config         `mapstructure:"api"`
	Checkers    checkers.Config    `mapstructure:"checkers"`
	Receiver    receiver.Config    `mapstructure:"receiver"`
}

func main() {
	defer logging.Init().Close()

	cfg := baseconfig.MustSetupAll(&Config{}, baseconfig.WithEnvPrefix("FASTAD_ALLINONE"))

	setupConfig(cfg)

	runCtx, shutdownCtx, cancel := stop.SetupCtx()
	defer cancel()

	g, gctx := errgroup.WithContext(runCtx)

	g.Go(func() error {
		if err := dataserviceImpl.Run(gctx, shutdownCtx, &cfg.DataService); err != nil {
			return fmt.Errorf("running dataservice: %w", err)
		}
		return nil
	})

	cfg.Checkers.DataService.Address = cfg.DataService.ListenAddress
	g.Go(func() error {
		if err := checkersImpl.Run(gctx, shutdownCtx, &cfg.Checkers); err != nil {
			return fmt.Errorf("running checkers: %w", err)
		}
		return nil
	})

	cfg.Scheduler.DataService.Address = cfg.DataService.ListenAddress
	g.Go(func() error {
		if err := schedulerImpl.Run(gctx, shutdownCtx, &cfg.Scheduler); err != nil {
			return fmt.Errorf("running scheduler: %w", err)
		}
		return nil
	})

	// Receiver and scoreboard depend on API, start it first.
	cfg.API.DataService.Address = cfg.DataService.ListenAddress
	cfg.API.ScoreboardAddress = cfg.Scoreboard.ListenAddress
	cfg.API.ReceiverAddress = cfg.Receiver.ListenAddress
	g.Go(func() error {
		if err := apiImpl.Run(gctx, shutdownCtx, &cfg.API); err != nil {
			return fmt.Errorf("running api: %w", err)
		}
		return nil
	})

	if err := apiwait.HTTP(gctx, cfg.API.ListenAddress); err != nil {
		zap.L().Error("waiting for API", zap.Error(err))
		cancel()
	}

	g.Go(func() error {
		if err := scoreboardImpl.Run(gctx, shutdownCtx, &cfg.Scoreboard); err != nil {
			return fmt.Errorf("running scoreboard: %w", err)
		}
		return nil
	})

	cfg.Receiver.DataService.Address = cfg.DataService.ListenAddress
	g.Go(func() error {
		if err := receiverImpl.Run(gctx, shutdownCtx, &cfg.Receiver); err != nil {
			return fmt.Errorf("running receiver: %w", err)
		}
		return nil
	})

	g.Go(func() error {
		metrics.RunServer(gctx, shutdownCtx, cfg.MetricsAddress)
		return nil
	})

	if err := g.Wait(); err != nil {
		zap.L().Fatal("all-in-one run failed", zap.Error(err))
	}
}

func setupConfig(cfg *Config) {
	if cfg.IntercomToken == "" {
		cfg.IntercomToken = uuid.NewString()
	}

	cfg.API.IntercomToken = cfg.IntercomToken
	cfg.Scoreboard.IntercomToken = cfg.IntercomToken
	cfg.Receiver.IntercomToken = cfg.IntercomToken

	cfg.DataService.Installation = fmt.Sprintf("%s/dataservice", cfg.UserAgent)
	cfg.Scheduler.Installation = fmt.Sprintf("%s/scheduler", cfg.UserAgent)
	cfg.Scoreboard.Installation = fmt.Sprintf("%s/scoreboard", cfg.UserAgent)
	cfg.API.Installation = fmt.Sprintf("%s/api", cfg.UserAgent)
	cfg.Checkers.Installation = fmt.Sprintf("%s/checkers", cfg.UserAgent)
	cfg.Receiver.Installation = fmt.Sprintf("%s/receiver", cfg.UserAgent)

	cfg.DataService.MetricsAddress = ""
	cfg.Scheduler.MetricsAddress = ""
	cfg.Scoreboard.MetricsAddress = ""
	cfg.API.MetricsAddress = ""
	cfg.Receiver.MetricsAddress = ""

	// Checkers are listening on a different port.
	cfg.Checkers.MetricsAddress = ""
}
