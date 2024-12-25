package checkers

import (
	"github.com/c4t-but-s4d/fastad/pkg/config"
)

type Config struct {
	Installation string `mapstructure:"installation" default:"checkers/worker"`

	MetricsAddress string `mapstructure:"metrics_address" default:":3006"`

	DataService config.DataService `mapstructure:"data_service"`
	Temporal    config.Temporal    `mapstructure:"temporal"`
	Postgres    config.Postgres    `mapstructure:"postgres"`
}
