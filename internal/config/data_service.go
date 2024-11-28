package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type DataService struct {
	Address string `mapstructure:"address"`
}

func SetDefaultDataServiceConfig(prefix string) {
	viper.SetDefault(fmt.Sprintf("%s.address", prefix), "127.0.0.1:1337")
}
