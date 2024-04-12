package dto

import "github.com/vuquang23/poseidon/internal/pkg/entity"

type GetSwapEventsReq struct {
	TxHash string `form:"txHash"`
}

type GetSwapEventsResp struct {
	Events []*entity.SwapEvent `json:"events"`
}
