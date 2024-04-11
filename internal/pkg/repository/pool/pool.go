package pool

import (
	"context"

	"github.com/vuquang23/poseidon/internal/pkg/entity"
)

type IPoolRepository interface {
	CreatePool(ctx context.Context, pool *entity.Pool) (uint64, error)
	GetPoolByAddress(ctx context.Context, address string) (*entity.Pool, error)

	CreateBlockCursors(ctx context.Context, cursors []*entity.BlockCursor) error
}
