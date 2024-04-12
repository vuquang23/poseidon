package entity

type SwapEvent struct {
	ID        uint64 `gorm:"primaryKey;autoIncrement" json:"id"`
	PoolID    uint64 `json:"poolId"`
	TxID      uint64 `json:"txId"`
	Amount0   string `json:"amount0"`
	Amount1   string `json:"amount1"`
	Price     string `json:"price"`
	CreatedAt uint64 `gorm:"autoCreateTime" json:"createdAt"`

	TxHash      string `gorm:"-" json:"txHash,omitempty"`
	PoolAddress string `gorm:"-" json:"poolAddress,omitempty"`
}
