package dataservice

import (
	"github.com/c4t-but-s4d/fastad/pkg/config"
)

type Config struct {
	Installation string `mapstructure:"installation" default:"dataservice"`

	IntercomToken string `mapstructure:"intercom_token"`

	ListenAddress  string `mapstructure:"listen_address" default:":8004"`
	MetricsAddress string `mapstructure:"metrics_address" default:":3004"`

	Postgres config.Postgres `mapstructure:"postgres"`
}
