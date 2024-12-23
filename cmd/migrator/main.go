package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/sirupsen/logrus"
	"github.com/uptrace/bun/migrate"
	"github.com/urfave/cli/v2"

	"github.com/c4t-but-s4d/fastad/cmd/migrator/migrations"
	"github.com/c4t-but-s4d/fastad/internal/config"
)

type Config struct {
	Postgres config.Postgres `mapstructure:"postgres"`
}

func main() {
	app := cli.NewApp()

	var cfg *Config
	var migrator *migrate.Migrator
	app.Before = func(_ *cli.Context) error {
		var err error
		if cfg, err = config.SetupAll(&Config{}, config.WithEnvPrefix("FASTAD_MIGRATOR")); err != nil {
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
					logrus.Info("there are no new migrations to run")
					return nil
				}

				logrus.Infof("migrated to %s", group)
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
					logrus.Info("there are no groups to roll back")
					return nil
				}

				logrus.Infof("rolled back %s", group)
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
				logrus.Infof("created migration %s (%s)", mf.Name, mf.Path)

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
					logrus.Infof("created migration %s (%s)", mf.Name, mf.Path)
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
				logrus.Infof("migrations: %s", ms)
				logrus.Infof("unapplied migrations: %s", ms.Unapplied())
				logrus.Infof("last migration group: %s", ms.LastGroup())

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
					logrus.Info("there are no new migrations to mark as applied")
					return nil
				}

				logrus.Infof("marked as applied %s", group)
				return nil
			},
		},
	}

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	if err := app.RunContext(ctx, os.Args); err != nil {
		logrus.WithError(err).Fatal("error running app")
	}
}
