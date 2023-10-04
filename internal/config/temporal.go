package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type Temporal struct {
	Address string `mapstructure:"address"`
}

func SetDefaultTemporalConfig(prefix string) {
	viper.SetDefault(fmt.Sprintf("%s.address", prefix), "temporal-server:7233")
}
