package tx

import (
	"context"

	"github.com/shopspring/decimal"

	"github.com/vuquang23/poseidon/internal/pkg/entity"
	"github.com/vuquang23/poseidon/internal/pkg/service/dto"
)

type ITxService interface {
	GetTxFeeUSDT(ctx context.Context, txHash string) (decimal.Decimal, error)
	GetSwapEventsByTxHash(ctx context.Context, txHash string) ([]*entity.SwapEvent, error)
	GetTxs(ctx context.Context, query dto.GetTxsQuery) ([]*entity.Tx, int64, error)
}
