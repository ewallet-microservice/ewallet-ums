package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/mhasnanr/ewallet-ums/bootstrap"
	"go.uber.org/zap"
)

func LoggerMiddleware(appLog bootstrap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		c.Next()

		duration := time.Since(start)
		method := c.Request.Method
		host := c.Request.Host
		urlPath := c.Request.URL.Path
		status := c.Writer.Status()

		if err := c.Errors.Last(); err != nil {
			appLog.Errorw(err.Error(),
				zap.Time("timestamp", start),
				zap.String("method", method),
				zap.String("host", host),
				zap.String("path", urlPath),
				zap.Int("status", status),
				zap.Int64("duration_ms", duration.Milliseconds()))

			return
		}

		appLog.Infow("request log",
			zap.Time("timestamp", start),
			zap.String("method", method),
			zap.String("host", host),
			zap.String("path", urlPath),
			zap.Int("status", status),
			zap.Int64("duration_ms", duration.Milliseconds()))
	}
}
