package handlers

import (
	"log"
	"platform/go-pkg/parser"
	"platform/go-pkg/response"
	"todo/types"

	"github.com/gin-gonic/gin"
)

func (h *TodoHandler) GetTodoDetailHandler(c *gin.Context) {
	log.Printf("hit service get todo detail with request: %v", c.Request)

	userID := c.GetString("user_id")

	todoID := c.Param("id")
	if todoID == "" {
		log.Println("missing todo id")
		response.BadRequest(c, "missing todo id", "todo id is required")
		return
	}
	
	log.Println("running controller")
	result, err := h.controller.GetTodoDetail(c.Request.Context(), &types.ReqGetTodoDetail{
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

	log.Println("success get todo detail")
	response.OK(c, "success", result)
}