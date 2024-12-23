package scoreboard

import (
	"github.com/c4t-but-s4d/fastad/internal/config"
)

type Config struct {
	UserAgent string `mapstructure:"user_agent" default:"scoreboard"`

	Postgres *config.Postgres `mapstructure:"postgres"`
}
