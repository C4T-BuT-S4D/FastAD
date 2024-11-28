package config

import (
	"fmt"
	"time"

	"github.com/spf13/viper"
)

type Postgres struct {
	Host      string `mapstructure:"host"`
	Port      int    `mapstructure:"port"`
	User      string `mapstructure:"user"`
	Password  string `mapstructure:"password"`
	Database  string `mapstructure:"database"`
	EnableSSL bool   `mapstructure:"enable_ssl"`

	MaxIdleConns    int           `mapstructure:"max_idle_conns"`
	MaxOpenConns    int           `mapstructure:"max_open_conns"`
	ConnMaxIdleTime time.Duration `mapstructure:"conn_max_idle_time"`
	ConnMaxLifetime time.Duration `mapstructure:"conn_max_lifetime"`
}

func SetDefaultPostgresConfig(prefix string) {
	viper.SetDefault(fmt.Sprintf("%s.host", prefix), "127.0.0.1")
	viper.SetDefault(fmt.Sprintf("%s.port", prefix), 5433)
	viper.SetDefault(fmt.Sprintf("%s.user", prefix), "fastad")
	viper.SetDefault(fmt.Sprintf("%s.password", prefix), "fastad")
	viper.SetDefault(fmt.Sprintf("%s.database", prefix), "fastad")
	viper.SetDefault(fmt.Sprintf("%s.enable_ssl", prefix), false)

	viper.SetDefault(fmt.Sprintf("%s.max_idle_conns", prefix), 10)
	viper.SetDefault(fmt.Sprintf("%s.max_open_conns", prefix), 10)
	viper.SetDefault(fmt.Sprintf("%s.conn_max_idle_time", prefix), 5*time.Minute)
	viper.SetDefault(fmt.Sprintf("%s.conn_max_lifetime", prefix), 10*time.Minute)
}
