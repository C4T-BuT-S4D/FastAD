package checkers

import (
	"github.com/c4t-but-s4d/fastad/pkg/config"
)

type Config struct {
	UserAgent string `mapstructure:"user_agent" default:"checkers/worker"`

	DataService config.DataService `mapstructure:"data_service"`
	Temporal    config.Temporal    `mapstructure:"temporal"`
	Postgres    config.Postgres    `mapstructure:"postgres"`
}
