package config

import (
	"fmt"

	"github.com/vlady-kotsev/lime-radio/shared/config"
	sharedconfig "github.com/vlady-kotsev/lime-radio/shared/config"

	"github.com/spf13/viper"
)

type Config struct {
	App   AppConfig                `mapstructure:"app"`
	Auth  AuthConfig               `mapstructure:"auth"`
	Radio RadioConfig              `mapstructure:"radio"`
	Event sharedconfig.EventConfig `mapstructure:"event"`
}

var _ sharedconfig.AuthConfiger = (*Config)(nil)
var _ sharedconfig.AppConfiger = (*Config)(nil)
var _ RadioConfiger = (*Config)(nil)
var _ sharedconfig.EventConfiger = (*Config)(nil)

type AppConfig struct {
	Port uint32 `mapstructure:"port"`
}

type AuthConfig struct {
	SharedSecret           config.Secret `mapstructure:"shared_secret"`
	AllowedOrigins         []string      `mapstructure:"allowed_origins"`
	TokenExpirationMinutes int64         `mapstructure:"expiration_minutes"`
}

type RadioConfig struct {
	SongsFolder string `mapstructure:"songs_folder"`
}

func Load() (*Config, error) {
	v := viper.New()
	v.SetConfigName("config")
	v.SetConfigType("toml")
	v.AddConfigPath(".")
	v.AddConfigPath("./config")

	if err := v.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("read config: %w", err)
	}

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("unmarshal config: %w", err)
	}

	return &cfg, nil
}

func (c *Config) GetTokenExpirationMinutes() int64 {
	return c.Auth.TokenExpirationMinutes
}

func (c *Config) GetSecretBytes() []byte {
	return c.Auth.SharedSecret.Bytes()
}

func (c *Config) GetAllowedOrigins() []string {
	return c.Auth.AllowedOrigins
}

func (c *Config) GetPort() uint32 {
	return c.App.Port
}

func (c *Config) GetSongFolder() string {
	return c.Radio.SongsFolder
}

func (c *Config) GetBrokerUrl() string {
	return c.Event.BrokerURL
}

func (c *Config) GetEventUsername() string {
	return c.Event.Username
}

func (c *Config) GetEventPassword() string {
	return string(c.Event.Password.Bytes())
}
