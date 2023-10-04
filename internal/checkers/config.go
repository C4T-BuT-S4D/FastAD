package checkers

import (
	"github.com/c4t-but-s4d/fastad/internal/config"
)

type DataService struct {
	Address string `mapstructure:"address"`
}

type Config struct {
	UserAgent string `mapstructure:"user_agent"`

	DataService DataService     `mapstructure:"data_service"`
	Temporal    config.Temporal `mapstructure:"temporal"`
}
