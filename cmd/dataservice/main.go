package main

import (
	"go.uber.org/zap"

	"github.com/c4t-but-s4d/fastad/cmd/dataservice/impl"
	"github.com/c4t-but-s4d/fastad/internal/dataservice"
	"github.com/c4t-but-s4d/fastad/pkg/baseconfig"
	"github.com/c4t-but-s4d/fastad/pkg/logging"
	"github.com/c4t-but-s4d/fastad/pkg/stop"
)

func main() {
	defer logging.Init().Close()

	cfg := baseconfig.MustSetupAll(&dataservice.Config{}, baseconfig.WithEnvPrefix("FASTAD_DATA_SERVICE"))

	runCtx, shutdownCtx, cancel := stop.SetupCtx()
	defer cancel()

	if err := impl.Run(runCtx, shutdownCtx, cfg); err != nil {
		zap.L().Fatal("error running server", zap.Error(err))
	}
}
