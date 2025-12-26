package config

import (
	"fmt"

	"github.com/vlady-kotsev/lime-radio/shared/config"

	"github.com/spf13/viper"
)

type Config struct {
	App     AppConfig     `mapstructure:"app"`
	Auth    AuthConfig    `mapstructure:"auth"`
	Payment PaymentConfig `mapstructure:"payment"`
}

var _ config.AuthConfiger = (*Config)(nil)
var _ config.AppConfiger = (*Config)(nil)
var _ PaymentConfiger = (*Config)(nil)

type AppConfig struct {
	Port uint32 `mapstructure:"port"`
}

type AuthConfig struct {
	SharedSecret           config.Secret `mapstructure:"shared_secret"`
	AllowedOrigins         []string      `mapstructure:"allowed_origins"`
	TokenExpirationMinutes int64         `mapstructure:"expiration_minutes"`
}

type PaymentConfig struct {
	Network         string `mapstructure:"network"`
	ReceiverAddress string `mapstructure:"receiver_address"`
	AmountLamports  uint64 `mapstructure:"amount_lamports"`
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

func (c *Config) GetAllowedOrigins() []string {
	return c.Auth.AllowedOrigins
}

func (c *Config) GetSecretBytes() []byte {
	return c.Auth.SharedSecret.Bytes()
}

func (c *Config) GetPort() uint32 {
	return c.App.Port
}

func (c *Config) GetNetwork() string {
	return c.Payment.Network
}

func (c *Config) GetReceiverAddress() string {
	return c.Payment.ReceiverAddress
}
func (c *Config) GetAmountLamports() uint64 {
	return c.Payment.AmountLamports
}
