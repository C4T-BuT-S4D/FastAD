package scheduler

import "github.com/c4t-but-s4d/fastad/internal/config"

type Config struct {
	UserAgent string `mapstructure:"user_agent" default:"scheduler"`

	Postgres    config.Postgres    `mapstructure:"postgres"`
	Temporal    config.Temporal    `mapstructure:"temporal"`
	DataService config.DataService `mapstructure:"data_service"`
}
