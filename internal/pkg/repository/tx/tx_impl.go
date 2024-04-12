package tx

import (
	"context"

	"github.com/pkg/errors"
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
			txIDByTxHash := map[string]uint64{}
			for _, tx := range txs {
				txIDByTxHash[tx.TxHash] = tx.ID
			}

			for _, e := range swapEvents {
				e.TxID = txIDByTxHash[e.TxHash]
			}

			if err := tx.Create(&swapEvents).Error; err != nil {
				logger.Error(ctx, err.Error())
				return err
			}
		}

		return nil
	})
}

func (r *TxRepository) GetTxsByPoolIDAndBlockRange(ctx context.Context, poolID uint64, fromBlock, toBlock uint64) ([]*entity.Tx, error) {
	var txs []*entity.Tx
	if err := r.db.
		Where("pool_id = ?", poolID).
		Where("block_number >= ? AND block_number <= ?", fromBlock, toBlock).
		Find(&txs).Error; err != nil {
		logger.Error(ctx, err.Error())
		return nil, err
	}

	return txs, nil
}

func (r *TxRepository) UpdateDataFinalizer(
	ctx context.Context,
	poolID, cursorID uint64,
	fromBlock, toBlock uint64,
	newTxs []*entity.Tx, newSwapEvents []*entity.SwapEvent,
) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&entity.BlockCursor{}).
			Where("id = ?", cursorID).
			Update("block_number", toBlock+1).Error; err != nil {
			logger.Error(ctx, err.Error())
			return err
		}

		if len(newTxs) == 0 {
			err := tx.Model(&entity.Tx{}).
				Where("pool_id = ?", poolID).
				Where("block_number >= ? AND block_number <= ?", fromBlock, toBlock).
				Update("is_finalized", true).Error
			if err != nil {
				logger.Error(ctx, err.Error())
				return err
			}

			return nil
		}

		err := tx.Where("pool_id = ?", poolID).
			Where("block_number >= ? AND block_number <= ?", fromBlock, toBlock).
			Delete(&entity.Tx{}).Error
		if err != nil {
			logger.Error(ctx, err.Error())
			return err
		}

		if err := tx.Create(&newTxs).Error; err != nil {
			logger.Error(ctx, err.Error())
			return err
		}

		if len(newSwapEvents) == 0 {
			return nil
		}

		txIDByTxHash := map[string]uint64{}
		for _, tx := range newTxs {
			txIDByTxHash[tx.TxHash] = tx.ID
		}

		for _, e := range newSwapEvents {
			e.TxID = txIDByTxHash[e.TxHash]
		}

		if err := tx.Create(&newSwapEvents).Error; err != nil {
			logger.Error(ctx, err.Error())
			return err
		}

		return nil
	})
}

func (r *TxRepository) GetTxByHash(ctx context.Context, hash string) (*entity.Tx, error) {
	var tx entity.Tx
	if err := r.db.Where("tx_hash = ?", hash).Take(&tx).Error; err != nil {
		logger.Error(ctx, err.Error())
		return nil, errors.Wrap(ErrTxNotFound, err.Error())
	}

	return &tx, nil
}

func (r *TxRepository) GetSwapEventsByTxHash(ctx context.Context, txHash string) ([]*entity.SwapEvent, error) {
	type event struct {
		ID        uint64
		PoolID    uint64
		TxID      uint64
		Amount0   string
		Amount1   string
		Price     string
		CreatedAt uint64
		Address   string
	}

	var events []*event

	err := r.db.Transaction(func(tx *gorm.DB) error {
		subQuery := tx.Model(&entity.Tx{}).Where("tx_hash = ?", txHash).Select("id")

		err := tx.Model(&entity.SwapEvent{}).
			Where("tx_id IN (?)", subQuery).
			Joins("LEFT JOIN pools ON swap_events.pool_id = pools.id").
			Select("swap_events.id, swap_events.pool_id, swap_events.tx_id, swap_events.amount0, swap_events.amount1, swap_events.price, swap_events.created_at, pools.address").
			Scan(&events).Error
		if err != nil {
			logger.Error(ctx, err.Error())
			return err
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	ret := make([]*entity.SwapEvent, 0, len(events))
	for _, e := range events {
		ret = append(ret, &entity.SwapEvent{
			ID:          e.ID,
			PoolID:      e.PoolID,
			TxID:        e.TxID,
			Amount0:     e.Amount0,
			Amount1:     e.Amount1,
			Price:       e.Price,
			CreatedAt:   e.CreatedAt,
			PoolAddress: e.Address,
		})
	}

	return ret, nil
}
