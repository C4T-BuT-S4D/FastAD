package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/urfave/cli/v2"
	"go.uber.org/zap"

	"github.com/c4t-but-s4d/fastad/cmd/fastad/cli/common"
	"github.com/c4t-but-s4d/fastad/cmd/fastad/cli/run"
	"github.com/c4t-but-s4d/fastad/cmd/fastad/cli/tokens"
	"github.com/c4t-but-s4d/fastad/pkg/logging"
)

func main() {
	defer logging.Init().Close()

	app := cli.NewApp()

	runCtx, runCancel := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT)
	defer runCancel()

	cc := &common.CommandContext{}

	app.Commands = []*cli.Command{
		run.NewCommand(cc),
		tokens.NewCommand(cc),
	}

	if err := app.RunContext(runCtx, os.Args); err != nil {
		zap.L().Fatal("error running app", zap.Error(err))
	}
}
