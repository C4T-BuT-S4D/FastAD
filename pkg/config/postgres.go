package config

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"

	"github.com/c4t-but-s4d/fastad/pkg/logging"
)

type Postgres struct {
	DSN       string `mapstructure:"dsn"`
	Host      string `mapstructure:"host" default:"127.0.0.1"`
	Port      int    `mapstructure:"port" default:"5433"`
	User      string `mapstructure:"user" default:"fastad"`
	Password  string `mapstructure:"password" default:"fastad"`
	Database  string `mapstructure:"database" default:"fastad"`
	EnableSSL bool   `mapstructure:"enable_ssl"`

	MaxIdleConns    int           `mapstructure:"max_idle_conns" default:"10"`
	MaxOpenConns    int           `mapstructure:"max_open_conns" default:"10"`
	ConnMaxIdleTime time.Duration `mapstructure:"conn_max_idle_time" default:"5m"`
	ConnMaxLifetime time.Duration `mapstructure:"conn_max_lifetime" default:"10m"`
}

func (p *Postgres) BunDB() *bun.DB {
	var opts []pgdriver.Option
	if p.DSN != "" {
		opts = append(opts, pgdriver.WithDSN(p.DSN))
	} else {
		opts = append(opts,
			pgdriver.WithAddr(fmt.Sprintf("%s:%d", p.Host, p.Port)),
			pgdriver.WithDatabase(p.Database),
			pgdriver.WithUser(p.User),
			pgdriver.WithPassword(p.Password),
			pgdriver.WithInsecure(!p.EnableSSL),
		)
	}

	pgConn := pgdriver.NewConnector(opts...)

	sqlDB := sql.OpenDB(pgConn)
	sqlDB.SetMaxIdleConns(p.MaxIdleConns)
	sqlDB.SetMaxOpenConns(p.MaxOpenConns)
	sqlDB.SetConnMaxIdleTime(p.ConnMaxIdleTime)
	sqlDB.SetConnMaxLifetime(p.ConnMaxLifetime)

	db := bun.NewDB(sqlDB, pgdialect.New())
	logging.AddBunQueryHook(db)

	return db
}
