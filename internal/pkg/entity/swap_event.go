package entity

import "github.com/shopspring/decimal"

type SwapEvent struct {
	ID        uint64          `gorm:"primaryKey;autoIncrement" json:"id"`
	PoolID    uint64          `json:"poolId"`
	TxID      uint64          `json:"txId"`
	Amount0   decimal.Decimal `json:"amount0"`
	Amount1   decimal.Decimal `json:"amount1"`
	Price     decimal.Decimal `json:"price"`
	CreatedAt uint64          `gorm:"autoCreateTime" json:"createdAt"`
}
