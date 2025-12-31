package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/vlady-kotsev/lime-radio/dj/internal/domain"
	"github.com/vlady-kotsev/lime-radio/dj/internal/service/payment"
	"github.com/vlady-kotsev/lime-radio/dj/internal/service/transaction"
)

const (
	PaymentRequestHeader  string = "PAYMENT-REQUIRED"
	PaymentResponseHeader string = "PAYMENT"
)

type PaymentMiddleware struct {
	paymentService     payment.PaymentServicer
	transactionService transaction.TransactionServicer
}

func NewPaymentMiddleware(paymentService payment.PaymentServicer, transactionService transaction.TransactionServicer) *PaymentMiddleware {
	return &PaymentMiddleware{
		paymentService:     paymentService,
		transactionService: transactionService,
	}
}

func (pm *PaymentMiddleware) ImposePayment() fiber.Handler {
	return func(c *fiber.Ctx) error {
		if c.Path() != "/request" {
			return c.Next()
		}
		paymentHeader := c.Get(PaymentResponseHeader)
		if paymentHeader == "" {
			requestPayload, err := pm.paymentService.ConstructPaymentPayload()
			if err != nil {
				return c.SendStatus(fiber.StatusInternalServerError)
			}
			c.Set(PaymentRequestHeader, requestPayload)
			return c.Status(fiber.StatusPaymentRequired).SendString("PAYMENT required")
		}

		paymentResponse, err := domain.PaymentResponseFromBase64(paymentHeader)
		if err != nil {
			return c.SendStatus(fiber.StatusBadRequest)
		}

		txHash := paymentResponse.TxHash
		isValid := pm.transactionService.ValidateTransaction(txHash)
		if !isValid {
			return c.SendStatus(fiber.StatusBadRequest)
		}

		return c.Next()
	}
}
