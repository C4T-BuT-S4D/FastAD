package checkers

type DataService struct {
	Address string `mapstructure:"address"`
}

type Temporal struct {
	Address string `mapstructure:"address"`
}

type Config struct {
	UserAgent string `mapstructure:"user_agent"`

	DataService DataService `mapstructure:"data_service"`
	Temporal    Temporal    `mapstructure:"temporal"`
}
