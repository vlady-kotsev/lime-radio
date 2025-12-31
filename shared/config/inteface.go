package config

type AuthConfiger interface {
	GetTokenExpirationMinutes() int64
	GetAllowedOrigins() []string
	GetSecretBytes() []byte
}

type AppConfiger interface {
	GetPort() uint32
}

type EventConfiger interface {
	GetBrokerUrl() string
	GetEventUsername() string
	GetEventPassword() string
}

type Configer interface {
	AuthConfiger
	AppConfiger
}
