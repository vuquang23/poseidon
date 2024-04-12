package api

import (
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/vuquang23/poseidon/internal/pkg/api/dto"
	"github.com/vuquang23/poseidon/internal/pkg/api/validator"
	svcdto "github.com/vuquang23/poseidon/internal/pkg/service/dto"
	"github.com/vuquang23/poseidon/internal/pkg/service/tx"
	"github.com/vuquang23/poseidon/pkg/logger"
)

func GetTxFeeUSDT(txSvc tx.ITxService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req dto.GetTxFeeReq
		if err := c.ShouldBindQuery(&req); err != nil {
			logger.Error(c, err.Error())
			RespondFailure(c, err)
			return
		}

		feeUSDT, err := txSvc.GetTxFeeUSDT(c, strings.ToLower(req.TxHash))
		if err != nil {
			RespondFailure(c, err)
			return
		}

		RespondSuccess(c, dto.GetTxFeeResp{FeeUSDT: feeUSDT})
	}
}

func GetSwapEvents(txSvc tx.ITxService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req dto.GetSwapEventsReq
		if err := c.ShouldBindQuery(&req); err != nil {
			logger.Error(c, err.Error())
			RespondFailure(c, err)
			return
		}

		events, err := txSvc.GetSwapEventsByTxHash(c, strings.ToLower(req.TxHash))
		if err != nil {
			RespondFailure(c, err)
			return
		}

		RespondSuccess(c, dto.GetSwapEventsResp{Events: events})
	}
}

func GetTxs(txSvc tx.ITxService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req dto.GetTxsReq
		if err := c.ShouldBindQuery(&req); err != nil {
			logger.Error(c, err.Error())
			RespondFailure(c, err)
			return
		}

		if err := validator.ValidateGetTxsReq(&req); err != nil {
			RespondFailure(c, err)
			return
		}

		poolAddress := strings.ToLower(req.PoolAddress)
		page := req.Page
		if page == 0 {
			page = 1
		}
		pageSize := req.PageSize
		if pageSize == 0 {
			pageSize = 50
		}

		query := svcdto.GetTxsQuery{
			PoolAddress: poolAddress,
			Page:        page,
			PageSize:    pageSize,
		}
		txs, total, err := txSvc.GetTxs(c, query)
		if err != nil {
			RespondFailure(c, err)
			return
		}

		RespondSuccess(c, dto.GetTxsResp{
			Txs:        txs,
			Pagination: dto.PaginationResp{Total: total},
		})
	}
}
