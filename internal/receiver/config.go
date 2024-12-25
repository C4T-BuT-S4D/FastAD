package receiver

import (
	"github.com/c4t-but-s4d/fastad/pkg/config"
)

type Config struct {
	Installation string `mapstructure:"installation" default:"receiver"`

	ListenAddress  string `mapstructure:"listen_address" default:":8002"`
	MetricsAddress string `mapstructure:"metrics_address" default:":3002"`

	Postgres    config.Postgres    `mapstructure:"postgres"`
	DataService config.DataService `mapstructure:"data_service"`
}
