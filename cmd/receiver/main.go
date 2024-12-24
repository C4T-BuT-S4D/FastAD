package main

import (
	"go.uber.org/zap"

	"github.com/c4t-but-s4d/fastad/cmd/receiver/impl"
	"github.com/c4t-but-s4d/fastad/internal/receiver"
	"github.com/c4t-but-s4d/fastad/pkg/baseconfig"
	"github.com/c4t-but-s4d/fastad/pkg/logging"
	"github.com/c4t-but-s4d/fastad/pkg/stop"
)

func main() {
	defer logging.Init().Close()

	cfg := baseconfig.MustSetupAll(&receiver.Config{}, baseconfig.WithEnvPrefix("FASTAD_RECEIVER"))

	runCtx, shutdownCtx, cancel := stop.SetupCtx()
	defer cancel()

	if err := impl.Run(runCtx, shutdownCtx, cfg); err != nil {
		zap.L().Fatal("receiver run failed", zap.Error(err))
	}
}
