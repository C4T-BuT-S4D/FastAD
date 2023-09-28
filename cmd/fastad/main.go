package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/c4t-but-s4d/fastad/cmd/fastad/cli/common"
	"github.com/c4t-but-s4d/fastad/cmd/fastad/cli/setup"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

func main() {
	app := cli.NewApp()

	runCtx, runCancel := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT)
	defer runCancel()

	cc := &common.CommandContext{}

	app.Commands = []*cli.Command{
		setup.NewSetupCommand(cc),
	}

	if err := app.RunContext(runCtx, os.Args); err != nil {
		logrus.Fatalf("error running app: %v", err)
	}
}
