package scoreboard

import (
	"time"

	"github.com/c4t-but-s4d/fastad/internal/config"
)

type Config struct {
	UserAgent string `mapstructure:"user_agent" default:"scoreboard"`

	Postgres config.Postgres `mapstructure:"postgres"`

	CheckInterval time.Duration `mapstructure:"check_interval" default:"1s"`
	BatchSize     int           `mapstructure:"batch_size" default:"1000"`
}
