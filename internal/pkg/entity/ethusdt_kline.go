package entity

import "github.com/shopspring/decimal"

// ETHUSDTKline "1m" interval kline.
type ETHUSDTKline struct {
	ID         uint64 `gorm:"primaryKey;autoIncrement"`
	OpenTime   uint64
	CloseTime  uint64
	OpenPrice  string
	HighPrice  string
	LowPrice   string
	ClosePrice string
	OHLC4      decimal.Decimal `gorm:"ohlc4"`
	CreatedAt  uint64          `gorm:"autoCreateTime"`
}
