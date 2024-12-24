package config

type DataService struct {
	Address string `mapstructure:"address" default:"127.0.0.1:1337"`
}
