package handlers

import (
	"github.com/c4t-but-s4d/fastad/pkg/config"
)

type Config struct {
	Installation string `mapstructure:"installation" default:"api"`

	ListenAddress  string `mapstructure:"listen_address" default:":8001"`
	MetricsAddress string `mapstructure:"metrics_address" default:":3001"`

	IntercomToken string `mapstructure:"intercom_token"`

	Postgres    config.Postgres    `mapstructure:"postgres"`
	DataService config.DataService `mapstructure:"data_service"`

	ScoreboardChannel string `mapstructure:"scoreboard_channel" default:"scoreboard"`

	ReceiverAddress string `mapstructure:"receiver_address" default:"localhost:8002"`
	SlacAddress     string `mapstructure:"slac_address" default:"localhost:8003"`
}
