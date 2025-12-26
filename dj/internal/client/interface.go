package client

import (
	"context"

	"github.com/gagliardetto/solana-go/rpc"
)

type SolanaClienter interface {
	GetTransaction(ctx context.Context, txHash string) (*rpc.GetTransactionResult, error)
}
