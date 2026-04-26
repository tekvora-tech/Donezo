package middleware

import (
	"platform/go-pkg/jwt"
	"platform/go-pkg/logger"
	"platform/go-pkg/response"
	"time"

	"github.com/gin-gonic/gin"
)

// AuthConfig konfigurasi auth middleware
type AuthConfig struct {
	SecretKey string
	Logger    *logger.Logger
}

// Auth middleware cek JWT token
func Auth(cfg AuthConfig) gin.HandlerFunc {
	// Inisialisasi JWT service sekali saat middleware dibuat
	jwtConfig := jwt.Config{
		Secret:          cfg.SecretKey,
		AccessTokenTTL:  15 * time.Minute,  // atau dari env
		RefreshTokenTTL: 7 * 24 * time.Hour,
		Issuer:          "donezo",
	}
	jwtService, _ := jwt.NewService(jwtConfig)

	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			response.Unauthorized(c, "token tidak ditemukan", "token tidak ditemukan")
			c.Abort()
			return
		}

		// Extract Bearer token
		token, err := jwtService.ExtractBearerToken(authHeader)
		if err != nil {
			response.Unauthorized(c, "format token tidak valid", err.Error())
			c.Abort()
			return
		}

		// Validasi JWT
		claims, err := jwtService.ValidateAccessToken(token)
		if err != nil {
			response.Unauthorized(c, "token tidak valid", err.Error())
			c.Abort()
			return
		}

		// Simpan data user di context
		c.Set("user_id", claims.UserID.String())
		c.Set("email", claims.Email)

		c.Next()
	}
}

// OptionalAuth auth opsional
func OptionalAuth(cfg AuthConfig) gin.HandlerFunc {
	jwtConfig := jwt.Config{
		Secret:          cfg.SecretKey,
		AccessTokenTTL:  15 * time.Minute,
		RefreshTokenTTL: 7 * 24 * time.Hour,
		Issuer:          "donezo",
	}
	jwtService, _ := jwt.NewService(jwtConfig)

	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.Next()
			return
		}

		token, err := jwtService.ExtractBearerToken(authHeader)
		if err != nil {
			c.Next()
			return
		}

		claims, err := jwtService.ValidateAccessToken(token)
		if err != nil {
			c.Next()
			return
		}

		c.Set("user_id", claims.UserID.String())
		c.Set("email", claims.Email)

		c.Next()
	}
}