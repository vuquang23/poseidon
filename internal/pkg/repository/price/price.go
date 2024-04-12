package price

import (
	"context"

	"github.com/vuquang23/poseidon/internal/pkg/entity"
)

type IPriceRepository interface {
	CreateKline(ctx context.Context, e *entity.ETHUSDTKline) error
	GetKlineByOpenTime(ctx context.Context, openTimeNsec int64) (*entity.ETHUSDTKline, error)
}
