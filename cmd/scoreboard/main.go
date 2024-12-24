package main

import (
	"github.com/c4t-but-s4d/fastad/internal/baseconfig"
	"github.com/c4t-but-s4d/fastad/internal/logging"
	"github.com/c4t-but-s4d/fastad/internal/scoreboard"
)

func main() {
	defer logging.Init().Close()

	_ = baseconfig.MustSetupAll(&scoreboard.Config{}, baseconfig.WithEnvPrefix("FASTAD_DATA_SERVICE"))
}
