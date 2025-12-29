package domain

import (
	"encoding/base64"
	"encoding/json"
)

type PaymentRequest struct {
	Network        string `json:"network"`
	PayTo          string `json:"payTo"`
	AmountLamports uint64 `json:"amount_lamports"`
	Description    string `json:"description,omitempty"`
}

func NewPaymentRequest(amountLamports uint64, network, payTo, description string) *PaymentRequest {
	return &PaymentRequest{
		AmountLamports: amountLamports,
		Network:        network,
		PayTo:          payTo,
		Description:    description,
	}
}

func (pr *PaymentRequest) ToBase64() (string, error) {
	b, err := json.Marshal(pr)
	if err != nil {
		return "", err
	}

	encoded := base64.StdEncoding.EncodeToString(b)
	return encoded, nil
}

type PaymentResponse struct {
	Network string `json:"network"`
	TxHash  string `json:"transaction_hash"`
}

func PaymentResponseFromBase64(encodedPaymentRequest string) (*PaymentResponse, error) {
	decodedBytes, err := base64.StdEncoding.DecodeString(encodedPaymentRequest)
	if err != nil {
		return nil, err
	}
	var paymentResponse PaymentResponse
	err = json.Unmarshal(decodedBytes, &paymentResponse)
	if err != nil {
		return nil, err
	}

	return &paymentResponse, nil
}
