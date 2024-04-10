package pool

import (
	"context"

	"github.com/vuquang23/poseidon/internal/pkg/service/dto"
)

type IPoolService interface {
	CreatePool(ctx context.Context, cmd *dto.CreatePoolCmd) (uint64, error)
}
