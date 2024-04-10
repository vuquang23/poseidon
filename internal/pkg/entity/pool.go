package entity

type Pool struct {
	ID             uint64 `gorm:"primaryKey;autoIncrement" json:"id"`
	Address        string `json:"address"`
	Token0         string `json:"token0"`
	Token0Decimals uint   `json:"token0Decimals"`
	Token1         string `json:"token1"`
	Token1Decimals uint   `json:"token1Decimals"`
	CreatedAt      uint64 `gorm:"autoCreateTime" json:"createdAt"`
}
