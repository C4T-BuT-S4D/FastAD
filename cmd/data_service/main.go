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
	"github.com/c4t-but-s4d/fastad/internal/services/teams"
	teamspb "github.com/c4t-but-s4d/fastad/pkg/proto/data/teams"
	"github.com/mitchellh/mapstructure"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type Postgres struct {
	Host      string `mapstructure:"host"`
	Port      int    `mapstructure:"port"`
	User      string `mapstructure:"user"`
	Password  string `mapstructure:"password"`
	Database  string `mapstructure:"database"`
	EnableSSL bool   `mapstructure:"enable_ssl"`
}

type Redis struct {
	Host string `mapstructure:"host"`
	Port int    `mapstructure:"port"`
	DB   int    `mapstructure:"db"`
}

type Config struct {
	ListenAddress string   `mapstructure:"listen_address"`
	Postgres      Postgres `mapstructure:"postgres"`
	Redis         Redis    `mapstructure:"redis"`
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

	redisClient := redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%d", cfg.Redis.Host, cfg.Redis.Port),
		DB:   cfg.Redis.DB,
	})

	teamsController := teams.NewController(db, redisClient)
	teamsService := teams.NewService(teamsController)

	server := grpc.NewServer()
	reflection.Register(server)
	teamspb.RegisterTeamsServiceServer(server, teamsService)

	runCtx, runCancel := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT)
	defer runCancel()

	if err := db.PingContext(runCtx); err != nil {
		logrus.Fatalf("error pinging postgres: %v", err)
	}
	if err := redisClient.Ping(runCtx).Err(); err != nil {
		logrus.Fatalf("error pinging redis: %v", err)
	}

	if err := teamsController.Migrate(runCtx); err != nil {
		logrus.Fatalf("error migrating teams: %v", err)
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
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))

	viper.SetDefault("listen_address", "127.0.0.1:11337")

	viper.SetDefault("postgres.host", "127.0.0.1")
	viper.SetDefault("postgres.port", 5432)
	viper.SetDefault("postgres.user", "local")
	viper.SetDefault("postgres.password", "local")
	viper.SetDefault("postgres.database", "local")

	viper.SetDefault("redis.host", "127.0.0.1")
	viper.SetDefault("redis.port", 6379)
	viper.SetDefault("redis.db", 0)

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