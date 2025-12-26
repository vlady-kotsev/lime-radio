package config

type AuthConfiger interface {
	GetTokenExpirationMinutes() int64
	GetAllowedOrigins() []string
	GetSecretBytes() []byte
}

type AppConfiger interface {
	GetPort() uint32
}

type Configer interface {
	AuthConfiger
	AppConfiger
}
