package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type Config struct {
	App  AppConfig  `mapstructure:"app"`
	Auth AuthConfig `mapstructure:"auth"`
}

type AppConfig struct {
	Port        uint32 `mapstructure:"port"`
	SongsFolder string `mapstructure:"songs_folder"`
}

type AuthConfig struct {
	SharedSecret   Secret   `mapstructure:"shared_secret"`
	AllowedOrigins []string `mapstructure:"allowed_origins"`
	TokenExpirationMinutes uint `mapstructure:"expiration_minutes"`
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
