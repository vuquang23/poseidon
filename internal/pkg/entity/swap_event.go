package entity

type SwapEvent struct {
	ID        uint64 `gorm:"primaryKey;autoIncrement" json:"id"`
	PoolID    uint64 `json:"poolId"`
	TxHash    string `json:"txHash"`
	Amount0   string `json:"amount0"`
	Amount1   string `json:"amount1"`
	Price     string `json:"price"`
	CreatedAt uint64 `gorm:"autoCreateTime" json:"createdAt"`
}
