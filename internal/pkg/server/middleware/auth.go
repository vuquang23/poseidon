package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

const HeaderXAPIKey = "X-API-KEY"

func NewAuthMiddleware(apiKey string) gin.HandlerFunc {
	return func(c *gin.Context) {
		k := c.GetHeader(HeaderXAPIKey)
		if k != apiKey {
			c.AbortWithStatus(http.StatusForbidden)
			return
		}

		c.Next()
	}
}
