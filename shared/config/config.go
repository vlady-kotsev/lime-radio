package config

type EventConfig struct {
	BrokerURL string `mapstructure:"broker_url"`
	Username  string `mapstructure:"username"`
	Password  Secret `mapstructure:"password"`
}
