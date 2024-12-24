package api

import (
	"github.com/c4t-but-s4d/fastad/pkg/config"
)

type Config struct {
	UserAgent string `mapstructure:"user_agent" default:"api"`

	ListenAddress string `mapstructure:"listen_address" default:"localhost:8001"`

	DataService config.DataService `mapstructure:"data_service"`

	ReceiverAddress   string `mapstructure:"receiver_address" default:"localhost:8002"`
	ScoreboardAddress string `mapstructure:"scoreboard_address" default:"localhost:8003"`
}
