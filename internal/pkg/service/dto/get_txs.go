package dto

type GetTxsQuery struct {
	PoolAddress string
	Page        uint
	PageSize    uint
}
