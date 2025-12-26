package transaction

type TransactionServicer interface {
	ValidateTransaction(txHash string) bool
}
