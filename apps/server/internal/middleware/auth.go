package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/price-comparison/server/internal/response"
)

func APIKeyAuth(apiKey string) gin.HandlerFunc {
	return func(c *gin.Context) {
		if apiKey == "" {
			c.Next()
			return
		}

		key := c.GetHeader("X-API-Key")
		if key == "" {
			authHeader := c.GetHeader("Authorization")
			if strings.HasPrefix(authHeader, "Bearer ") {
				key = strings.TrimSpace(strings.TrimPrefix(authHeader, "Bearer "))
			}
		}

		if key != apiKey {
			response.Error(c, http.StatusUnauthorized, response.ErrUnauthorized, "invalid api key")
			c.Abort()
			return
		}

		c.Next()
	}
}
