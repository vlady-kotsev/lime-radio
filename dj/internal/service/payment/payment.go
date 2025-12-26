package payment

import (
	"github.com/vlady-kotsev/lime-radio/dj/internal/config"
	"github.com/vlady-kotsev/lime-radio/dj/internal/domain"
	"go.uber.org/zap"
)

type PaymentService struct {
	logger *zap.Logger
	config config.PaymentConfiger
}

var _ PaymentServicer = (*PaymentService)(nil)

func NewPaymentService(logger *zap.Logger, config config.PaymentConfiger) *PaymentService {
	return &PaymentService{logger: logger, config: config}
}

func (ps *PaymentService) ConstructPaymentPayload(description string) (string, error) {
	paymentRequest := domain.NewPaymentRequest(ps.config.GetAmountLamports(), ps.config.GetNetwork(), ps.config.GetReceiverAddress(), description)
	encodedPaymentRequest, err := paymentRequest.ToBase64()
	if err != nil {
		return "", err
	}
	return encodedPaymentRequest, nil
}
