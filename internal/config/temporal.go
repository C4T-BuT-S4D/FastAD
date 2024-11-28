package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type Temporal struct {
	Address string `mapstructure:"address"`
}

func SetDefaultTemporalConfig(prefix string) {
	viper.SetDefault(fmt.Sprintf("%s.address", prefix), "127.0.0.1:7233")
}
