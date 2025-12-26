package transaction

import (
	"context"
	"time"

	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
	"github.com/vlady-kotsev/lime-radio/dj/internal/client"
	"github.com/vlady-kotsev/lime-radio/dj/internal/config"
	"go.uber.org/zap"
)

type TransactionService struct {
	logger *zap.Logger
	config config.PaymentConfiger
	client client.SolanaClienter
}

var _ TransactionServicer = (*TransactionService)(nil)

func NewTransactionService(logger *zap.Logger, config config.PaymentConfiger, client client.SolanaClienter) *TransactionService {
	return &TransactionService{
		logger: logger,
		config: config,
		client: client,
	}
}

func (ts *TransactionService) ValidateTransaction(txHash string) bool {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	txResult, err := ts.client.GetTransaction(ctx, txHash)
	if err != nil {
		ts.logger.Error("Failed to get transaction", zap.Error(err), zap.String("txHash", txHash))
		return false
	}

	if txResult == nil {
		ts.logger.Warn("Transaction result is nil", zap.String("txHash", txHash))
		return false
	}

	if txResult.Meta != nil && txResult.Meta.Err != nil {
		ts.logger.Warn("Transaction failed on-chain", zap.Any("error", txResult.Meta.Err), zap.String("txHash", txHash))
		return false
	}

	tx, err := txResult.Transaction.GetTransaction()
	if err != nil {
		ts.logger.Error("Failed to parse transaction", zap.Error(err), zap.String("txHash", txHash))
		return false
	}

	if ts.validatePaymentDetails(txResult, tx) {
		ts.logger.Info("Transaction validation successful", zap.String("txHash", txHash))
		return true
	}

	ts.logger.Warn("Transaction validation failed", zap.String("txHash", txHash))
	return false
}

func (ts *TransactionService) validatePaymentDetails(txResult *rpc.GetTransactionResult, tx *solana.Transaction) bool {
	expectedAmount := ts.config.GetAmountLamports()
	receiverAddress := ts.config.GetReceiverAddress()

	receiverPubkey, err := solana.PublicKeyFromBase58(receiverAddress)
	if err != nil {
		ts.logger.Error("Invalid receiver address in config", zap.Error(err), zap.String("receiverAddress", receiverAddress))
		return false
	}

	if txResult.Meta != nil && txResult.Meta.PreBalances != nil && txResult.Meta.PostBalances != nil {
		for i, postBalance := range txResult.Meta.PostBalances {
			if i >= len(txResult.Meta.PreBalances) {
				continue
			}

			balanceChange := int64(postBalance) - int64(txResult.Meta.PreBalances[i])

			if balanceChange >= int64(expectedAmount) {
				if i < len(tx.Message.AccountKeys) && tx.Message.AccountKeys[i].Equals(receiverPubkey) {

					return true
				}
			}
		}
	}

	return false
}
