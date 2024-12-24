package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/urfave/cli/v2"
	"go.uber.org/zap"

	"github.com/c4t-but-s4d/fastad/cmd/fastad/cli/common"
	"github.com/c4t-but-s4d/fastad/cmd/fastad/cli/setup"
	"github.com/c4t-but-s4d/fastad/internal/logging"
)

func main() {
	defer logging.Init().Close()

	// TODO: rewrite with cobra.
	app := cli.NewApp()

	runCtx, runCancel := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT)
	defer runCancel()

	cc := &common.CommandContext{}

	app.Commands = []*cli.Command{
		setup.NewSetupCommand(cc),
	}

	if err := app.RunContext(runCtx, os.Args); err != nil {
		zap.L().With(zap.Error(err)).Fatal("error running app")
	}
}
