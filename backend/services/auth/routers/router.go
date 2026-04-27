package routers

import (
	"auth/handlers"

	"github.com/gin-gonic/gin"
)

func RegisterAuthRoutes(r *gin.RouterGroup, h *handlers.AuthHandler, auth gin.HandlerFunc) {
	g := r.Group("/auth")
	{	
		// Public
		g.POST("/v1/register", h.RegisterHandler)
		g.POST("/v1/login", h.LoginHandler)
		g.POST("/v1/refresh", h.RefreshTokenHandler)

		// Protected
		g.POST("/v1/logout", auth, h.LogoutHandler)
	}
}