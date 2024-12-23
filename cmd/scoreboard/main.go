package main

import (
	"github.com/c4t-but-s4d/fastad/internal/config"
	"github.com/c4t-but-s4d/fastad/internal/logging"
	"github.com/c4t-but-s4d/fastad/internal/scoreboard"
)

func main() {
	_ = config.MustSetupAll(&scoreboard.Config{}, config.WithEnvPrefix("FASTAD_DATA_SERVICE"))

	logging.Init()
}
