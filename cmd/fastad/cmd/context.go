package cmd

import (
	"context"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type Context struct {
	Context context.Context //nolint:containedctx // cobra pattern: context passed through command tree
	Viper   *viper.Viper
}

type CommandFactory func(cc *Context) *cobra.Command

func NewContext(v *viper.Viper) *Context {
	return &Context{
		Viper: v,
	}
}
