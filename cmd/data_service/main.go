package main

import (
	"context"
	"database/sql"
	"fmt"
	"net"
	"os/signal"
	"strings"
	"syscall"

	"github.com/c4t-but-s4d/fastad/internal/config"
	"github.com/c4t-but-s4d/fastad/internal/logging"
	"github.com/c4t-but-s4d/fastad/internal/services/gamestate"
	"github.com/c4t-but-s4d/fastad/internal/services/services"
	"github.com/c4t-but-s4d/fastad/internal/services/teams"
	"github.com/c4t-but-s4d/fastad/internal/version"
	"github.com/c4t-but-s4d/fastad/pkg/grpctools"
	gspb "github.com/c4t-but-s4d/fastad/pkg/proto/data/game_state"
	servicespb "github.com/c4t-but-s4d/fastad/pkg/proto/data/services"
	teamspb "github.com/c4t-but-s4d/fastad/pkg/proto/data/teams"
	"github.com/mitchellh/mapstructure"
	"github.com/sirupsen/logrus"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
)

type Config struct {
	ListenAddress string          `mapstructure:"listen_address"`
	Postgres      config.Postgres `mapstructure:"postgres"`
}

func main() {
	cfg, err := setupConfig()
	if err != nil {
		logrus.Fatalf("error setting up config: %v", err)
	}

	logging.Init()

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

	versionController := version.NewController(db)

	teamsController := teams.NewController(db, versionController)
	teamsService := teams.NewService(teamsController)

	servicesController := services.NewController(db, versionController)
	servicesService := services.NewService(servicesController)

	gameStateController := gamestate.NewController(db, versionController)
	gameStateService := gamestate.NewService(gameStateController)

	server := grpctools.NewServer()
	teamspb.RegisterTeamsServiceServer(server, teamsService)
	servicespb.RegisterServicesServiceServer(server, servicesService)
	gspb.RegisterGameStateServiceServer(server, gameStateService)

	runCtx, runCancel := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT)
	defer runCancel()

	if err := db.PingContext(runCtx); err != nil {
		logrus.Fatalf("error pinging postgres: %v", err)
	}

	if err := versionController.Migrate(runCtx); err != nil {
		logrus.Fatalf("error migrating versions: %v", err)
	}

	if err := teamsController.Migrate(runCtx); err != nil {
		logrus.Fatalf("error migrating teams: %v", err)
	}

	if err := servicesController.Migrate(runCtx); err != nil {
		logrus.Fatalf("error migrating services: %v", err)
	}

	if err := gameStateController.Migrate(runCtx); err != nil {
		logrus.Fatalf("error migrating game state: %v", err)
	}

	logrus.Infof("starting server on %s", cfg.ListenAddress)
	lis, err := net.Listen("tcp", cfg.ListenAddress)
	if err != nil {
		logrus.Fatalf("error listening on %s: %v", cfg.ListenAddress, err)
	}

	go func() {
		<-runCtx.Done()
		server.GracefulStop()
	}()

	if err := server.Serve(lis); err != nil {
		logrus.Fatalf("error serving: %v", err)
	}
}

func setupConfig() (*Config, error) {
	pflag.BoolP("debug", "v", false, "Enable verbose logging")
	pflag.Parse()
	if err := viper.BindPFlags(pflag.CommandLine); err != nil {
		return nil, fmt.Errorf("binding pflags: %w", err)
	}
	viper.SetEnvPrefix("FASTAD_DATA_SERVICE")
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_", ".", "_"))

	viper.SetDefault("listen_address", "127.0.0.1:1337")

	config.SetDefaultPostgresConfig("postgres")

	cfg := new(Config)
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
