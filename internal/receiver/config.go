package receiver

import (
	"github.com/c4t-but-s4d/fastad/pkg/config"
)

type Config struct {
	UserAgent string `mapstructure:"user_agent" default:"receiver"`

	ListenAddress string `mapstructure:"listen_address" default:":8002"`

	Postgres    config.Postgres    `mapstructure:"postgres"`
	DataService config.DataService `mapstructure:"data_service"`
}
