package entity

import "github.com/shopspring/decimal"

type ETHUSDTKline struct {
	ID         uint64 `gorm:"primaryKey;autoIncrement"`
	OpenTime   uint64
	CloseTime  uint64
	OpenPrice  decimal.Decimal
	HighPrice  decimal.Decimal
	LowPrice   decimal.Decimal
	ClosePrice decimal.Decimal
	CreatedAt  uint64 `gorm:"autoCreateTime"`
}
