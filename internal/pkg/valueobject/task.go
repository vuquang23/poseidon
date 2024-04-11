package valueobject

type TaskType string

const (
	TaskTypeHandlePoolCreated = "handle_pool_created"
	TaskTypeGetETHUSDTKline   = "get_ethusdt_kline"
	TaskTypeScanTxs           = "scan_txs"
)

type TaskHandlePoolCreatedPayload struct {
	PoolAddress string `json:"poolAddress"`
}

type TaskScanTxsPayload struct {
	PoolID         uint64 `json:"poolId"`
	PoolAddress    string `json:"poolAddress"`
	Token0Decimals uint   `json:"token0Decimals"`
	Token1Decimals uint   `json:"token1Decimals"`
}

type TaskGetETHUSDTKlinePayload struct {
	Time uint64 `json:"time"`
}
