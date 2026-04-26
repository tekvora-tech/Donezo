package middleware

import (
	"fmt"
	"runtime/debug"

	"platform/go-pkg/logger"
	"platform/go-pkg/response"

	"github.com/gin-gonic/gin"
)

// Recovery tangkap panic dan return error 500
func Recovery(log *logger.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				// Log error + stack trace
				log.WithFields(logger.Fields{
					"error":      fmt.Sprintf("%v", err),
					"stacktrace": string(debug.Stack()),
					"path":       c.Request.URL.Path,
					"method":     c.Request.Method,
				}).Error("panic recovered")

				// Response ke client (jangan expose stack trace)
				response.InternalServerError(c, "terjadi kesalahan pada server", "terjadi kesalahan pada server")
				c.Abort()
			}
		}()

		c.Next()
	}
}