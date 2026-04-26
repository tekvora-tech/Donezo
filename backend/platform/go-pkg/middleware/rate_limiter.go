package middleware

import (
	"time"

	"platform/go-pkg/response"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

// RateLimiter limit request per IP
func RateLimiter(redisClient *redis.Client, limit int, window time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()
		key := "rate_limit:" + ip

		ctx := c.Request.Context()

		// Cek count
		count, err := redisClient.Get(ctx, key).Int()
		if err != nil && err != redis.Nil {
			c.Next()
			return
		}

		if count >= limit {
			response.TooManyRequests(c, "terlalu banyak request, coba lagi nanti", "terlalu banyak request, coba lagi nanti")
			c.Abort()
			return
		}

		// Increment count
		pipe := redisClient.Pipeline()
		pipe.Incr(ctx, key)
		pipe.Expire(ctx, key, window)
		_, _ = pipe.Exec(ctx)

		c.Next()
	}
}