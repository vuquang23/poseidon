package entity

import (
	"gorm.io/datatypes"

	"github.com/vuquang23/poseidon/internal/pkg/valueobject"
)

type BlockCursor struct {
	ID          uint64 `gorm:"primaryKey;autoIncrement"`
	PoolID      uint64
	Type        valueobject.BlockCursorType
	BlockNumber uint64
	Extra       datatypes.JSON
	CreatedAt   uint64 `gorm:"autoCreateTime"`
	UpdatedAt   uint64 `gorm:"autoUpdateTime"`
}
