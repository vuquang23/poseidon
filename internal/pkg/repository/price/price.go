package price

import (
	"context"

	"github.com/vuquang23/poseidon/internal/pkg/entity"
)

type IPriceRepository interface {
	CreateKline(ctx context.Context, e *entity.ETHUSDTKline) error
}
