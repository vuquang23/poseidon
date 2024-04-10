package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	poolsvc "github.com/vuquang23/poseidon/internal/pkg/service/pool"

	"github.com/vuquang23/poseidon/internal/pkg/config"
	"github.com/vuquang23/poseidon/internal/pkg/server/middleware"
)

func RegisterRoutes(
	conf *config.Config,
	server *gin.Engine,
	poolSvc poolsvc.IPoolService,
) {
	router := server.Group("/api/v1")

	/// health
	router.GET("health/live", func(c *gin.Context) { c.AbortWithStatusJSON(http.StatusOK, "OK") })
	router.GET("health/ready", func(c *gin.Context) { c.AbortWithStatusJSON(http.StatusOK, "OK") })

	/// pool
	router.POST("pools", middleware.NewAuthMiddleware(conf.Common.APIKey), CreatePool(poolSvc))
}
