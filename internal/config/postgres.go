package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type Postgres struct {
	Host      string `mapstructure:"host"`
	Port      int    `mapstructure:"port"`
	User      string `mapstructure:"user"`
	Password  string `mapstructure:"password"`
	Database  string `mapstructure:"database"`
	EnableSSL bool   `mapstructure:"enable_ssl"`
}

func SetDefaultPostgresConfig(prefix string) {
	viper.SetDefault(fmt.Sprintf("%s.host", prefix), "127.0.0.1")
	viper.SetDefault(fmt.Sprintf("%s.port", prefix), 5432)
	viper.SetDefault(fmt.Sprintf("%s.user", prefix), "local")
	viper.SetDefault(fmt.Sprintf("%s.password", prefix), "local")
	viper.SetDefault(fmt.Sprintf("%s.database", prefix), "local")
}
