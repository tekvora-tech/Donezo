package routers

import (
	"users/handlers"

	"github.com/gin-gonic/gin"
)

func RegisterUserRoutes(r *gin.RouterGroup, h *handlers.UserHandler, auth gin.HandlerFunc) {
	g := r.Group("/users")
	{	
		// Protected
		g.GET("/v1/profile", auth, h.GetProfileHandler)
		g.PUT("/v1/profile", auth, h.UpdateProfileHandler)
	}
}