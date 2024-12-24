package baseconfig

import (
	"fmt"
	"strings"

	"github.com/creasty/defaults"
	"github.com/mitchellh/mapstructure"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"go.uber.org/zap"
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

func SetupAll[T any](cfg *T, opts ...SetupOption) (*T, error) {
	setupConfig := GetSetupConfig(opts...)

	v, err := NewViper(setupConfig.envPrefix)
	if err != nil {
		return nil, fmt.Errorf("creating viper: %w", err)
	}

	return Setup(v, cfg)
}

func MustSetupAll[T any](cfg *T, opts ...SetupOption) *T {
	t, err := SetupAll[T](cfg, opts...)
	if err != nil {
		zap.L().Fatal("error setting up config", zap.Error(err))
	}

	return t
}
