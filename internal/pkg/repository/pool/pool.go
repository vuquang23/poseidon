package pool

import (
	"context"

	"github.com/vuquang23/poseidon/internal/pkg/entity"
)

type IPoolRepository interface {
	CreatePool(ctx context.Context, pool *entity.Pool) (uint64, error)
}
