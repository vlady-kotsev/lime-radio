package payment

type PaymentServicer interface {
	ConstructPaymentPayload(description string) (string, error)
}
