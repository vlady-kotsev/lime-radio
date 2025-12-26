package config

type PaymentConfiger interface {
	GetNetwork() string
	GetReceiverAddress() string
	GetAmountLamports() uint64
}
