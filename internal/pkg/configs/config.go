package configs

type (
	Config struct {
		App   App   `mapstructure:"app"`
		Chain Chain `mapstructure:"chain"`
	}

	Chain struct {
	}

	App struct {
		MetricAddr   string   `mapstructure:"metric_addr"`
		Env          string   `mapstructure:"env"`
		LogLevel     string   `mapstructure:"log_level"`
		ChannelTypes []string `mapstructure:"channel_types"`
	}
)

func NewConfig() *Config {
	return &Config{}
}
