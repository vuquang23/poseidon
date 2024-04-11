package tx

import (
	"context"

	"github.com/vuquang23/poseidon/internal/pkg/entity"
	"github.com/vuquang23/poseidon/internal/pkg/valueobject"
)

type ITxRepository interface {
	// UpdateDataScanner updates block cursor, persists new scanned txs and swap events for pool
	UpdateDataScanner(ctx context.Context, blockcursorID uint64, newBlockNbr uint64, txs []*entity.Tx, swapEvents []*entity.SwapEvent) error

	CreateBlockCursors(ctx context.Context, cursors []*entity.BlockCursor) error
	GetCursorByPoolIDAndType(ctx context.Context, poolID uint64, t valueobject.BlockCursorType) (*entity.BlockCursor, error)
}
