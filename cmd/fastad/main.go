package main

import (
	"context"
	"os/signal"
	"syscall"

	"github.com/spf13/cobra"
	"go.uber.org/zap"

	"github.com/c4t-but-s4d/fastad/cmd/fastad/cmd"
	"github.com/c4t-but-s4d/fastad/pkg/logging"
	"github.com/c4t-but-s4d/fastad/pkg/viperext"
)

func main() {
	cobra.EnableTraverseRunHooks = true

	defer logging.Init().Close()

	v, err := viperext.NewViper("FASTAD")
	if err != nil {
		zap.L().Fatal("creating viper", zap.Error(err))
	}

	cc := cmd.NewContext(v)

	rootCmd := cmd.NewRootCmd(cc)
	rootCmd.AddCommand(cmd.NewInitCmd(cc))
	rootCmd.AddCommand(cmd.NewRunCmd(cc))
	rootCmd.AddCommand(cmd.NewResetCmd(cc))
	rootCmd.AddCommand(cmd.NewTokensCmd(cc))

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT)
	defer cancel()

	if err := rootCmd.ExecuteContext(ctx); err != nil {
		zap.L().Fatal("running app", zap.Error(err))
	}
}
