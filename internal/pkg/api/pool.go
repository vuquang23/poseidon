package api

import (
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/vuquang23/poseidon/internal/pkg/api/dto"
	"github.com/vuquang23/poseidon/internal/pkg/api/validator"
	svcdto "github.com/vuquang23/poseidon/internal/pkg/service/dto"
	"github.com/vuquang23/poseidon/internal/pkg/service/pool"
	"github.com/vuquang23/poseidon/pkg/logger"
)

func CreatePool(poolsvc pool.IPoolService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req dto.CreatePoolReq
		if err := c.ShouldBindJSON(&req); err != nil {
			logger.Error(c, err.Error())
			RespondFailure(c, err)
			return
		}

		if err := validator.ValidateCreatePoolReq(&req); err != nil {
			RespondFailure(c, err)
			return
		}

		command := svcdto.CreatePoolCmd{
			Address:        strings.ToLower(req.Address),
			StartBlock:     req.StartBlock,
			Token0:         strings.ToLower(req.Token0),
			Token0Decimals: req.Token0Decimals,
			Token1:         strings.ToLower(req.Token1),
			Token1Decimals: req.Token1Decimals,
		}
		poolID, err := poolsvc.CreatePool(c, &command)
		if err != nil {
			RespondFailure(c, err)
			return
		}

		RespondSuccess(c, dto.CreatePoolResp{PoolID: poolID})
	}
}
