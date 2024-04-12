package tx

import (
	"context"

	"github.com/shopspring/decimal"
)

type ITxService interface {
	GetTxFeeUSDT(ctx context.Context, txHash string) (decimal.Decimal, error)
}
