package dto

import "github.com/shopspring/decimal"

type GetTxFeeReq struct {
	TxHash string `form:"txHash"`
}

type GetTxFeeResp struct {
	FeeUSDT decimal.Decimal `json:"feeUsdt"`
}
