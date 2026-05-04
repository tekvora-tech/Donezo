package routers

import (
	"todo/handlers"

	"github.com/gin-gonic/gin"
)

func RegisterTodoRoutes(r *gin.RouterGroup, h *handlers.TodoHandler, auth gin.HandlerFunc) {
	g := r.Group("/todos")
	{	
		// Protected
		g.POST("/v1/", auth, h.CreateTodoHandler)
		g.GET("/v1/", auth, h.ListTodosHandler)
		g.GET("/v1/:id", auth, h.GetTodoDetailHandler)
		g.PUT("/v1/:id", auth, h.UpdateTodoHandler)
		g.DELETE("/v1/:id", auth, h.DeleteTodoHandler)
	}
}