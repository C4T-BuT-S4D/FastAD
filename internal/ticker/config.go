package ticker

import (
	"github.com/c4t-but-s4d/fastad/internal/config"
)

type Config struct {
	Postgres config.Postgres `mapstructure:"postgres"`
	Temporal config.Temporal `mapstructure:"temporal"`
}
