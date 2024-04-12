package tx

import (
	"context"

	"github.com/shopspring/decimal"

	"github.com/vuquang23/poseidon/internal/pkg/entity"
)

type ITxService interface {
	GetTxFeeUSDT(ctx context.Context, txHash string) (decimal.Decimal, error)
	GetSwapEventsByTxHash(ctx context.Context, txHash string) ([]*entity.SwapEvent, error)
}
