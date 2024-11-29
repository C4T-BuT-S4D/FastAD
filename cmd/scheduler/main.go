package main

import (
	"context"
	"database/sql"
	"fmt"
	"os/signal"
	"syscall"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/spf13/pflag"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
	"go.temporal.io/sdk/client"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/c4t-but-s4d/fastad/internal/clients/gamestate"
	"github.com/c4t-but-s4d/fastad/internal/config"
	"github.com/c4t-but-s4d/fastad/internal/logging"
	"github.com/c4t-but-s4d/fastad/internal/scheduler"
	"github.com/c4t-but-s4d/fastad/pkg/grpctools"
	gspb "github.com/c4t-but-s4d/fastad/pkg/proto/data/game_state"
	"github.com/c4t-but-s4d/fastad/pkg/util"
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
				"component": "scheduler",
			}),
		),
	})
	if err != nil {
		logrus.Fatalf("dialing temporal: %v", err)
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

	gameStateClient := gamestate.NewClient(gspb.NewGameStateServiceClient(dataServiceConn))

	t := scheduler.New(time.Second*10, temporalClient, gameStateClient, db)

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	if err := t.MigrateDB(ctx); err != nil {
		logrus.Fatalf("migrating database: %v", err)
	}

	if err := t.Run(ctx); err != nil {
		logrus.Fatalf("running scheduler: %v", err)
	}
}

func setupConfig() (*scheduler.Config, error) {
	pflag.BoolP("debug", "v", false, "Enable verbose logging")
	pflag.Parse()

	cfg := util.Must[*scheduler.Config]("setup config")(
		config.SetupAll[*scheduler.Config]("FASTAD_DATA_SERVICE"),
	)

	return cfg, nil
}
