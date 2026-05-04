package handlers

import (
	"log"
	"platform/go-pkg/parser"
	"platform/go-pkg/response"
	"todo/types"

	"github.com/gin-gonic/gin"
)

func (h *TodoHandler) DeleteTodoHandler(c *gin.Context) {
	log.Printf("hit service delete todo with request: %v", c.Request)

	userID := c.GetString("user_id")

	todoID := c.Param("id")
	if todoID == "" {
		log.Println("missing todo id")
		response.BadRequest(c, "missing todo id", "todo id is required")
		return
	}

	log.Println("running controller")
	_, err := h.controller.DeleteTodo(c.Request.Context(), &types.ReqDeleteTodo{
		UserID: userID,
		TodoID: todoID,
	})

	if err != nil {
		log.Printf("error on controller: %v", err)
		if parser.IsPostgresError(err) {
			httpError := parser.ParsePostgresError(err)
			response.Error(c, httpError.Code, httpError.Message, httpError.Err)
			return
		}

		msg, code, detail := parser.ParseError(err)
		response.Error(c, code, msg, detail)
		return
	}

	log.Println("success delete todo")
	response.NoContent(c)
}