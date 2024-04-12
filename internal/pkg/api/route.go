package api

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/vuquang23/poseidon/internal/pkg/config"
	"github.com/vuquang23/poseidon/internal/pkg/server/middleware"
	poolsvc "github.com/vuquang23/poseidon/internal/pkg/service/pool"
	txsvc "github.com/vuquang23/poseidon/internal/pkg/service/tx"
)

func RegisterRoutes(
	conf *config.Config,
	server *gin.Engine,
	poolSvc poolsvc.IPoolService,
	txSvc txsvc.ITxService,
) {
	router := server.Group("/api/v1")

	// health
	router.GET("health/live", func(c *gin.Context) { c.AbortWithStatusJSON(http.StatusOK, "OK") })
	router.GET("health/ready", func(c *gin.Context) { c.AbortWithStatusJSON(http.StatusOK, "OK") })

	// pool
	router.POST("pools", middleware.NewAuthMiddleware(conf.Common.APIKey), CreatePool(poolSvc))

	// tx
	router.GET("txs/fee-usdt", GetTxFeeUSDT(txSvc))
	router.GET("txs/swap-events", GetSwapEvents(txSvc))
}
