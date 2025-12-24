package config

type AuthConfiger interface {
	GetTokenExpirationMinutes() int64
	GetAllowedOrigins() []string
}
