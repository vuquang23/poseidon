package entity

import (
	"github.com/shopspring/decimal"
	"gorm.io/datatypes"
)

type Tx struct {
	ID          uint64          `gorm:"primaryKey;autoIncrement" json:"id"`
	PoolID      uint64          `json:"poolId"`
	TxHash      string          `json:"txHash"`
	BlockNumber uint64          `json:"blockNumber"`
	BlockTime   uint64          `json:"blockTime"`
	Gas         decimal.Decimal `json:"gas"`
	Receipt     datatypes.JSON  `json:"-"`
	IsFinalized bool            `json:"isFinalized"`
	CreatedAt   uint64          `gorm:"autoCreateTime" json:"createdAt"`
}