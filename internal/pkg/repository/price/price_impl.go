package price

import (
	"context"

	"gorm.io/gorm"

	"github.com/vuquang23/poseidon/internal/pkg/entity"
	"github.com/vuquang23/poseidon/pkg/logger"
)

type PriceRepository struct {
	db *gorm.DB
}

func New(db *gorm.DB) *PriceRepository {
	return &PriceRepository{
		db: db,
	}
}

func (r *PriceRepository) CreateKline(ctx context.Context, e *entity.ETHUSDTKline) error {
	if err := r.db.Create(e).Error; err != nil {
		logger.Error(ctx, err.Error())
		return err
	}

	return nil
}
