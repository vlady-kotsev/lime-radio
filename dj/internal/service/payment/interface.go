package payment

type PaymentServicer interface {
	ConstructPaymentPayload() (string, error)
}
