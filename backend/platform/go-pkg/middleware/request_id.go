package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// RequestID tambah unique ID tiap request (untuk tracing)
func RequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := c.GetHeader("X-Request-ID")
		if requestID == "" {
			requestID = uuid.New().String()
		}

		// Simpan di context Gin
		c.Set("request_id", requestID)

		// Kirim balik ke client
		c.Writer.Header().Set("X-Request-ID", requestID)

		c.Next()
	}
}