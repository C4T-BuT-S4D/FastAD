package scoreboard

import (
	"time"

	"github.com/c4t-but-s4d/fastad/pkg/config"
)

type Config struct {
	Installation string `mapstructure:"installation" default:"scoreboard"`

	ListenAddress  string `mapstructure:"listen_address" default:":8003"`
	MetricsAddress string `mapstructure:"metrics_address" default:":3003"`

	Channel       string `mapstructure:"channel" default:"scoreboard"`
	IntercomToken string `mapstructure:"intercom_token"`

	Postgres         config.Postgres         `mapstructure:"postgres"`
	CentrifugeClient config.CentrifugeClient `mapstructure:"centrifuge_client"`

	CheckInterval            time.Duration `mapstructure:"check_interval" default:"1s"`
	BatchSize                int           `mapstructure:"batch_size" default:"1000"`
	ExecutionTXCreateTimeout time.Duration `mapstructure:"execution_tx_create_timeout" default:"5m"`
	ProcessorStateThreshold  time.Duration `mapstructure:"processor_state_threshold" default:"1m"`
	TimeCorrectionThreshold  time.Duration `mapstructure:"time_correction_threshold" default:"1m"`
}
