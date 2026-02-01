package viperext

import (
	"fmt"
	"strings"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

func NewViper(envPrefix string) (*viper.Viper, error) {
	v := viper.NewWithOptions(viper.ExperimentalBindStruct())

	if err := v.BindPFlags(pflag.CommandLine); err != nil {
		return nil, fmt.Errorf("binding pflags: %w", err)
	}
	v.SetEnvPrefix(envPrefix)
	v.AutomaticEnv()
	v.SetEnvKeyReplacer(
		strings.NewReplacer(
			"-", "_",
			".", "_",
		),
	)

	return v, nil
}
