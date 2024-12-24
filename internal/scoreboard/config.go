package scoreboard

import (
	"time"

	"github.com/c4t-but-s4d/fastad/internal/config"
)

type Config struct {
	UserAgent string `mapstructure:"user_agent" default:"scoreboard"`

	Postgres config.Postgres `mapstructure:"postgres"`

	CheckInterval            time.Duration `mapstructure:"check_interval" default:"1s"`
	BatchSize                int           `mapstructure:"batch_size" default:"1000"`
	ExecutionTXCreateTimeout time.Duration `mapstructure:"execution_tx_create_timeout" default:"5m"`
	ProcessorStateThreshold  time.Duration `mapstructure:"processor_state_threshold" default:"1m"`
	TimeCorrectionThreshold  time.Duration `mapstructure:"time_correction_threshold" default:"1m"`
}
