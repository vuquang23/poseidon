package pool

import (
	"context"

	"github.com/pkg/errors"

	"github.com/vuquang23/poseidon/internal/pkg/entity"
	poolrepo "github.com/vuquang23/poseidon/internal/pkg/repository/pool"
	"github.com/vuquang23/poseidon/internal/pkg/service/dto"
)

type PoolService struct {
	poolRepo poolrepo.IPoolRepository
}

func New(poolRepo poolrepo.IPoolRepository) *PoolService {
	return &PoolService{
		poolRepo: poolRepo,
	}
}

func (s *PoolService) CreatePool(ctx context.Context, cmd *dto.CreatePoolCmd) (uint64, error) {
	pool := entity.Pool{
		Address:        cmd.Address,
		Token0:         cmd.Token0,
		Token0Decimals: cmd.Token0Decimals,
		Token1:         cmd.Token1,
		Token1Decimals: cmd.Token1Decimals,
	}

	poolID, err := s.poolRepo.CreatePool(ctx, &pool)
	if err != nil {
		return 0, errors.Wrap(ErrCreatePool, err.Error())
	}

	return poolID, nil
}
