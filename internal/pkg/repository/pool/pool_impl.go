package pool

import (
	"context"

	"gorm.io/gorm"

	"github.com/vuquang23/poseidon/internal/pkg/entity"
	"github.com/vuquang23/poseidon/internal/pkg/valueobject"
	"github.com/vuquang23/poseidon/pkg/asynq"
	"github.com/vuquang23/poseidon/pkg/logger"
)

type PoolRepository struct {
	db          *gorm.DB
	asynqClient asynq.IAsynqClient
}

func New(db *gorm.DB, asynqClient asynq.IAsynqClient) *PoolRepository {
	return &PoolRepository{
		db:          db,
		asynqClient: asynqClient,
	}
}

func (r *PoolRepository) CreatePool(ctx context.Context, pool *entity.Pool) (uint64, error) {
	err := r.db.Transaction(func(tx *gorm.DB) error {
		if err := r.db.Create(pool).Error; err != nil {
			logger.Error(ctx, err.Error())
			return err
		}

		return r.asynqClient.EnqueueTask(
			ctx,
			string(valueobject.TaskTypeHandlePoolCreated),
			"", "",
			valueobject.TaskHandlePoolCreatedPayload{PoolAddress: pool.Address},
			-1,
		)
	})

	if err != nil {
		return 0, err
	}

	return pool.ID, nil
}
