package dataservice

import (
	"github.com/c4t-but-s4d/fastad/pkg/config"
)

type Config struct {
	ListenAddress string          `mapstructure:"listen_address" default:":8004"`
	Postgres      config.Postgres `mapstructure:"postgres"`
}
