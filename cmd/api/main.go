package main

import (
	"go.uber.org/zap"

	"github.com/c4t-but-s4d/fastad/cmd/api/impl"
	"github.com/c4t-but-s4d/fastad/internal/api"
	"github.com/c4t-but-s4d/fastad/pkg/baseconfig"
	"github.com/c4t-but-s4d/fastad/pkg/logging"
	"github.com/c4t-but-s4d/fastad/pkg/stop"
)

func main() {
	defer logging.Init().Close()

	cfg := baseconfig.MustSetupAll(&api.Config{}, baseconfig.WithEnvPrefix("API"))

	runCtx, shutdownCtx, cancel := stop.SetupCtx()
	defer cancel()

	if err := impl.Run(runCtx, shutdownCtx, cfg); err != nil {
		zap.L().Fatal("api run failed", zap.Error(err))
	}
}
