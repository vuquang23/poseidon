package tx

import (
	"context"

	"gorm.io/gorm"

	"github.com/vuquang23/poseidon/internal/pkg/entity"
	"github.com/vuquang23/poseidon/internal/pkg/valueobject"
	"github.com/vuquang23/poseidon/pkg/logger"
)

type TxRepository struct {
	db *gorm.DB
}

func New(db *gorm.DB) *TxRepository {
	return &TxRepository{
		db: db,
	}
}

func (r *TxRepository) CreateBlockCursors(ctx context.Context, cursors []*entity.BlockCursor) error {
	if err := r.db.Create(cursors).Error; err != nil {
		logger.Error(ctx, err.Error())
		return err
	}

	return nil
}

func (r *TxRepository) GetCursorByPoolIDAndType(ctx context.Context, poolID uint64, t valueobject.BlockCursorType) (*entity.BlockCursor, error) {
	var cursor entity.BlockCursor
	if err := r.db.Where("pool_id = ? AND type = ?", poolID, t).Take(&cursor).Error; err != nil {
		logger.Error(ctx, err.Error())
		return nil, err
	}

	return &cursor, nil
}

func (r *TxRepository) UpdateDataScanner(ctx context.Context, blockcursorID uint64, newBlockNbr uint64, txs []*entity.Tx, swapEvents []*entity.SwapEvent) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&entity.BlockCursor{}).
			Where("id = ?", blockcursorID).
			Update("block_number", newBlockNbr).Error; err != nil {
			logger.Error(ctx, err.Error())
			return err
		}

		if len(txs) != 0 {
			if err := tx.Create(&txs).Error; err != nil {
				logger.Error(ctx, err.Error())
				return err
			}
		}

		if len(swapEvents) != 0 {
			if err := tx.Create(&swapEvents).Error; err != nil {
				logger.Error(ctx, err.Error())
				return err
			}
		}

		return nil
	})
}
