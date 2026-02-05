package middleware

import (
	"log/slog"
	"time"

	"github.com/gin-gonic/gin"
)

func Logging(logger *slog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		c.Next()

		latency := time.Since(start)
		status := c.Writer.Status()

		logger.Info("request completed",
			"method", c.Request.Method,
			"path", c.Request.URL.Path,
			"status", status,
			"latency_ms", latency.Milliseconds(),
			"client_ip", c.ClientIP(),
			"user_agent", c.Request.UserAgent(),
			"request_id", c.GetString("request_id"),
		)
	}
}
