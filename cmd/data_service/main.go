package main

import (
	"context"
	"database/sql"
	"fmt"
	"net"
	"os/signal"
	"strings"
	"syscall"

	"github.com/c4t-but-s4d/fastad/internal/logging"
	"github.com/c4t-but-s4d/fastad/internal/services/services"
	"github.com/c4t-but-s4d/fastad/internal/services/teams"
	"github.com/c4t-but-s4d/fastad/internal/version"
	"github.com/c4t-but-s4d/fastad/pkg/grpctools"
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

type Postgres struct {
	Host      string `mapstructure:"host"`
	Port      int    `mapstructure:"port"`
	User      string `mapstructure:"user"`
	Password  string `mapstructure:"password"`
	Database  string `mapstructure:"database"`
	EnableSSL bool   `mapstructure:"enable_ssl"`
}

type Config struct {
	ListenAddress string   `mapstructure:"listen_address"`
	Postgres      Postgres `mapstructure:"postgres"`
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
	db := bun.NewDB(sqlDB, pgdialect.New())

	versionController := version.NewController(db)

	teamsController := teams.NewController(db, versionController)
	teamsService := teams.NewService(teamsController)

	servicesController := services.NewController(db, versionController)
	servicesService := services.NewService(servicesController)

	server := grpctools.NewServer()
	teamspb.RegisterTeamsServiceServer(server, teamsService)
	servicespb.RegisterServicesServiceServer(server, servicesService)

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

	viper.SetDefault("listen_address", "127.0.0.1:11337")

	viper.SetDefault("postgres.host", "127.0.0.1")
	viper.SetDefault("postgres.port", 5432)
	viper.SetDefault("postgres.user", "local")
	viper.SetDefault("postgres.password", "local")
	viper.SetDefault("postgres.database", "local")

	viper.SetDefault("redis.host", "127.0.0.1")
	viper.SetDefault("redis.port", 6379)
	viper.SetDefault("redis.db", 0)

	logrus.Infof("config: %+v", viper.AllSettings())

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
