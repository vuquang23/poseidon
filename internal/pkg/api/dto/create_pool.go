package dto

type CreatePoolReq struct {
	Address        string `json:"address"`
	Token0         string `json:"token0"`
	Token0Decimals uint   `json:"token0Decimals"`
	Token1         string `json:"token1"`
	Token1Decimals uint   `json:"token1Decimals"`
}

type CreatePoolResp struct {
	PoolID uint64 `json:"poolId"`
}
