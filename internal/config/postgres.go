package config

import (
	"time"
)

type Postgres struct {
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
