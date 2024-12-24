package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/uptrace/bun/migrate"
	"github.com/urfave/cli/v2"
	"go.uber.org/zap"

	"github.com/c4t-but-s4d/fastad/cmd/migrator/migrations"
	"github.com/c4t-but-s4d/fastad/internal/baseconfig"
	"github.com/c4t-but-s4d/fastad/internal/config"
	"github.com/c4t-but-s4d/fastad/internal/logging"
)

type Config struct {
	Postgres config.Postgres `mapstructure:"postgres"`
}

func main() {
	defer logging.Init().Close()

	app := cli.NewApp()

	var cfg *Config
	var migrator *migrate.Migrator
	app.Before = func(_ *cli.Context) error {
		var err error
		if cfg, err = baseconfig.SetupAll(&Config{}, baseconfig.WithEnvPrefix("FASTAD_MIGRATOR")); err != nil {
			return fmt.Errorf("setting up config: %w", err)
		}
		migrator = migrate.NewMigrator(cfg.Postgres.BunDB(), migrations.Migrations)
		return nil
	}

	app.Commands = []*cli.Command{
		{
			Name:  "init",
			Usage: "create migration tables",
			Action: func(c *cli.Context) error {
				return migrator.Init(c.Context)
			},
		},
		{
			Name:  "migrate",
			Usage: "migrate database",
			Action: func(c *cli.Context) error {
				group, err := migrator.Migrate(c.Context)
				if err != nil {
					return fmt.Errorf("migrating: %w", err)
				}

				if group.ID == 0 {
					zap.L().Info("there are no new migrations to run")
					return nil
				}

				zap.S().Infof("migrated to %s", group)
				return nil
			},
		},
		{
			Name:  "rollback",
			Usage: "rollback the last migration group",
			Action: func(c *cli.Context) error {
				group, err := migrator.Rollback(c.Context)
				if err != nil {
					return fmt.Errorf("rolling back: %w", err)
				}

				if group.ID == 0 {
					zap.L().Info("there are no groups to roll back")
					return nil
				}

				zap.S().Infof("rolled back %s", group)
				return nil
			},
		},
		{
			Name:  "create_go",
			Usage: "create Go migration",
			Action: func(c *cli.Context) error {
				name := strings.Join(c.Args().Slice(), "_")
				mf, err := migrator.CreateGoMigration(c.Context, name)
				if err != nil {
					return fmt.Errorf("creating Go migration: %w", err)
				}
				zap.S().Infof("created migration %s (%s)", mf.Name, mf.Path)

				return nil
			},
		},
		{
			Name:  "create_sql",
			Usage: "create up and down SQL migrations",
			Action: func(c *cli.Context) error {
				name := strings.Join(c.Args().Slice(), "_")
				files, err := migrator.CreateSQLMigrations(c.Context, name)
				if err != nil {
					return fmt.Errorf("creating SQL migrations: %w", err)
				}

				for _, mf := range files {
					zap.S().Infof("created migration %s (%s)", mf.Name, mf.Path)
				}

				return nil
			},
		},
		{
			Name:  "status",
			Usage: "print migrations status",
			Action: func(c *cli.Context) error {
				ms, err := migrator.MigrationsWithStatus(c.Context)
				if err != nil {
					return fmt.Errorf("checking migrations status: %w", err)
				}
				zap.S().Infof("migrations: %s", ms)
				zap.S().Infof("unapplied migrations: %s", ms.Unapplied())
				zap.S().Infof("last migration group: %s", ms.LastGroup())

				return nil
			},
		},
		{
			Name:  "mark_applied",
			Usage: "mark migrations as applied without actually running them",
			Action: func(c *cli.Context) error {
				group, err := migrator.Migrate(c.Context, migrate.WithNopMigration())
				if err != nil {
					return fmt.Errorf("migrating: %w", err)
				}

				if group.ID == 0 {
					zap.L().Info("there are no new migrations to mark as applied")
					return nil
				}

				zap.S().Infof("marked as applied %s", group)
				return nil
			},
		},
	}

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	if err := app.RunContext(ctx, os.Args); err != nil {
		zap.L().With(zap.Error(err)).Fatal("error running app")
	}
}
