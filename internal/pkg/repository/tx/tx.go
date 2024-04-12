package tx

import (
	"context"

	"github.com/vuquang23/poseidon/internal/pkg/entity"
	"github.com/vuquang23/poseidon/internal/pkg/valueobject"
)

type ITxRepository interface {
	// UpdateDataScanner updates block cursor, persists new scanned txs and swap events for pool
	UpdateDataScanner(ctx context.Context, blockcursorID uint64, newBlockNbr uint64, txs []*entity.Tx, swapEvents []*entity.SwapEvent) error

	// UpdateDataFinalizer updates block cursor, and
	// if reorg occured, staled txs and swap events are replaced by finalized ones.
	// Else, it updates column `is_finalized` of existing txs to be "true".
	UpdateDataFinalizer(ctx context.Context, poolID, cursorID uint64, fromBlock, toBlock uint64, newTxs []*entity.Tx, newSwapEvents []*entity.SwapEvent) error

	CreateBlockCursors(ctx context.Context, cursors []*entity.BlockCursor) error
	GetCursorByPoolIDAndType(ctx context.Context, poolID uint64, t valueobject.BlockCursorType) (*entity.BlockCursor, error)

	GetTxsByPoolIDAndBlockRange(ctx context.Context, poolID uint64, fromBlock, toBlock uint64) ([]*entity.Tx, error)
	GetTxByHash(ctx context.Context, hash string) (*entity.Tx, error)

	GetSwapEventsByTxHash(ctx context.Context, txHash string) ([]*entity.SwapEvent, error)
}
