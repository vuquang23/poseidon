package server

import (
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/requestid"
	"github.com/gin-gonic/gin"

	"github.com/vuquang23/poseidon/internal/pkg/server/middleware"
	requestidpkg "github.com/vuquang23/poseidon/internal/pkg/util/requestid"
	"github.com/vuquang23/poseidon/pkg/logger"
)

func GinEngine(config Config, logCfg logger.Config, logBackend logger.LoggerBackend) *gin.Engine {
	gin.SetMode(config.Mode)

	middlewares := []gin.HandlerFunc{
		requestid.New(
			requestid.WithCustomHeaderStrKey(requestidpkg.HeaderKeyRequestID),
			requestid.WithHandler(func(c *gin.Context, requestID string) {
				c.Request = c.Request.WithContext(requestidpkg.SetRequestIDToContext(c.Request.Context(), requestID))
			}),
		),
		middleware.NewLoggerMiddleware(logCfg, logBackend),
		gin.Recovery(),
	}

	engine := gin.New()
	engine.Use(middlewares...)

	setCORS(engine)

	return engine
}

func setCORS(engine *gin.Engine) {
	corsConfig := cors.DefaultConfig()
	corsConfig.AddAllowMethods(http.MethodOptions)
	corsConfig.AllowAllOrigins = true
	engine.Use(cors.New(corsConfig))
}
