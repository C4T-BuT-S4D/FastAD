package config

import (
	"fmt"
	"strings"

	"github.com/creasty/defaults"
	"github.com/mitchellh/mapstructure"
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
	v.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))

	return v, nil
}

func Setup[T any](v *viper.Viper, cfg T) (t T, err error) {
	if err := defaults.Set(cfg); err != nil {
		return t, fmt.Errorf("setting defaults: %w", err)
	}

	if err := v.Unmarshal(
		cfg,
		viper.DecodeHook(
			mapstructure.ComposeDecodeHookFunc(
				mapstructure.TextUnmarshallerHookFunc(),
				mapstructure.StringToTimeDurationHookFunc(),
			),
		),
	); err != nil {
		return t, fmt.Errorf("unmarshaling config: %w", err)
	}

	return cfg, nil
}

func SetupAll[T any](envPrefix string) (t T, err error) {
	v, err := NewViper(envPrefix)
	if err != nil {
		return t, fmt.Errorf("creating viper: %w", err)
	}

	return Setup(v, t)
}
