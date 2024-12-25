package scheduler

import (
	"github.com/c4t-but-s4d/fastad/pkg/config"
)

type Config struct {
	Installation string `mapstructure:"installation" default:"scheduler"`

	MetricsAddress string `mapstructure:"metrics_address" default:":3005"`

	Postgres    config.Postgres    `mapstructure:"postgres"`
	Temporal    config.Temporal    `mapstructure:"temporal"`
	DataService config.DataService `mapstructure:"data_service"`
}
