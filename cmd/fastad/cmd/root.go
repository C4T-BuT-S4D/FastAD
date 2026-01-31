package cmd

import (
	"github.com/spf13/cobra"

	"github.com/c4t-but-s4d/fastad/pkg/viperext"
)

func NewRootCmd(cc *Context) *cobra.Command {
	rootCmd := &cobra.Command{
		Use:              "fastad",
		Short:            "FastAD - Blazing fast Attack & Defence CTF platform",
		SilenceUsage:     true,
		SilenceErrors:    true,
		TraverseChildren: true,
		PersistentPreRunE: func(cmd *cobra.Command, _ []string) error {
			viperext.RunBindCommandFlags(cc.Viper, cmd)
			cc.Context = cmd.Context()
			return nil
		},
	}

	return rootCmd
}
