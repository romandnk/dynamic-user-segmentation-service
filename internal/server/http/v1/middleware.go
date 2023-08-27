package v1

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"time"
)

func (h *Handler) loggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		c.Next()

		duration := time.Since(start)

		h.logger.Info("Request info HTTP",
			zap.String("client ip", c.RemoteIP()),
			zap.String("method", c.Request.Method),
			zap.String("method path", c.FullPath()),
			zap.String("HTTP version", c.Request.Proto),
			zap.Int("status code", c.Writer.Status()),
			zap.String("processing time", duration.String()),
		)
	}
}
