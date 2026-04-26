package middleware

import (
	"time"

	"platform/go-pkg/logger"

	"github.com/gin-gonic/gin"
)

// Logger log tiap HTTP request
func Logger(log *logger.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		method := c.Request.Method

		c.Next()

		// Data setelah request selesai
		latency := time.Since(start)
		status := c.Writer.Status()
		clientIP := c.ClientIP()
		requestID, _ := c.Get("request_id")

		// Log sesuai level
		entry := log.WithFields(logger.Fields{
			"request_id": requestID,
			"method":     method,
			"path":       path,
			"status":     status,
			"latency":    latency.Milliseconds(),
			"client_ip":  clientIP,
		})

		if status >= 500 {
			entry.Error("server error")
		} else if status >= 400 {
			entry.Warn("client error")
		} else {
			entry.Info("request completed")
		}
	}
}