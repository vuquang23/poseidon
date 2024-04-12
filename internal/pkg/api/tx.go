package api

import (
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/vuquang23/poseidon/internal/pkg/api/dto"
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
