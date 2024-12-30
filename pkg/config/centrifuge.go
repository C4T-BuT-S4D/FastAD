package config

type CentrifugeClient struct {
	Address string `mapstructure:"address" default:"ws://127.0.0.1:8001/centrifuge/websocket"`
}
