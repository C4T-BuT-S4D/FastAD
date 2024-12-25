package api

import (
	"github.com/c4t-but-s4d/fastad/pkg/config"
)

type Config struct {
	Installation string `mapstructure:"installation" default:"api"`

	ListenAddress  string `mapstructure:"listen_address" default:":8001"`
	MetricsAddress string `mapstructure:"metrics_address" default:":3001"`

	DataService config.DataService `mapstructure:"data_service"`

	ReceiverAddress   string `mapstructure:"receiver_address" default:"localhost:8002"`
	ScoreboardAddress string `mapstructure:"scoreboard_address" default:"localhost:8003"`
}
