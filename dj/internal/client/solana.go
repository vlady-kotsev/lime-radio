package client

import (
	"context"
	"fmt"

	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
	"go.uber.org/zap"
)

type SolanaClient struct {
	client *rpc.Client
	logger *zap.Logger
	rpcURL string
}

var _ SolanaClienter = (*SolanaClient)(nil)

func NewSolanaClient(logger *zap.Logger) *SolanaClient {

	rpcURL := "https://devnet.helius-rpc.com/?api-key=2e770d32-a64b-458b-8cb8-06e208fba7f6"

	return &SolanaClient{
		client: rpc.New(rpcURL),
		logger: logger,
		rpcURL: rpcURL,
	}
}

func (sc *SolanaClient) GetTransaction(ctx context.Context, txHash string) (*rpc.GetTransactionResult, error) {
	sc.logger.Info("Fetching transaction", zap.String("txHash", txHash))

	signature, err := solana.SignatureFromBase58(txHash)
	if err != nil {
		sc.logger.Error("Invalid transaction signature", zap.Error(err), zap.String("txHash", txHash))
		return nil, fmt.Errorf("invalid transaction signature: %w", err)
	}

	result, err := sc.client.GetTransaction(
		ctx,
		signature,
		&rpc.GetTransactionOpts{
			Commitment: rpc.CommitmentConfirmed,
		},
	)
	if err != nil {
		sc.logger.Error("Failed to get transaction", zap.Error(err), zap.String("txHash", txHash))
		return nil, fmt.Errorf("failed to get transaction: %w", err)
	}

	if result == nil {
		sc.logger.Warn("Transaction not found", zap.String("txHash", txHash))
		return nil, fmt.Errorf("transaction not found: %s", txHash)
	}

	sc.logger.Info("Transaction fetched successfully",
		zap.String("txHash", txHash),
		zap.Any("slot", result.Slot),
	)

	return result, nil
}
