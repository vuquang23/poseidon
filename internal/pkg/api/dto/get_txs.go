package dto

import "github.com/vuquang23/poseidon/internal/pkg/entity"

type GetTxsReq struct {
	PoolAddress string `form:"poolAddress"`
	PaginationReq
}

type GetTxsResp struct {
	Txs        []*entity.Tx   `json:"txs"`
	Pagination PaginationResp `json:"pagination"`
}
