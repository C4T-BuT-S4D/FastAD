package baseconfig

type SetupConfig struct {
	envPrefix string
}

func GetSetupConfig(opts ...SetupOption) *SetupConfig {
	cfg := &SetupConfig{}
	for _, opt := range opts {
		opt(cfg)
	}
	return cfg
}

type SetupOption func(cfg *SetupConfig)

func WithEnvPrefix(envPrefix string) SetupOption {
	return func(cfg *SetupConfig) {
		cfg.envPrefix = envPrefix
	}
}
