package handlers

import (
	"log"
	"platform/go-pkg/parser"
	"platform/go-pkg/response"
	"todo/types"

	"github.com/gin-gonic/gin"
)

func (h *TodoHandler) UpdateTodoHandler(c *gin.Context) {
	log.Printf("hit service update todo with request: %v", c.Request)

	userID := c.GetString("user_id")

	todoID := c.Param("id")
	if todoID == "" {
		log.Println("missing todo id")
		response.BadRequest(c, "missing todo id", "todo id is required")
		return
	}

	var parsedBody types.DTOUpdateTodo
	if err := c.ShouldBindJSON(&parsedBody); err != nil {
		log.Printf("failed unmarshal request body: %v", err)
		response.BadRequest(c, "failed unmarshal request body", err.Error())
		return
	}

	if err := h.validator.Validate(parsedBody); err != nil {
		log.Printf("failed format request body: %v", err)
		response.BadRequest(c, "failed format request body", err.Error())
		return
	}

	log.Println("running controller")
	result, err := h.controller.UpdateTodo(c.Request.Context(), &types.ReqUpdateTodo{
		UserID: userID,
		TodoID: todoID,
		Title: parsedBody.Title,
		Description: parsedBody.Description,
		Status: parsedBody.Status,
		Priority: parsedBody.Priority,
		DueDate: parsedBody.DueDate,
		TagIDs: parsedBody.TagIDs,
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

	log.Println("success update todo")
	response.OK(c, "success", result)
}