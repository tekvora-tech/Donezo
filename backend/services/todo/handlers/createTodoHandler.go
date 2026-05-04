package handlers

import (
	"log"
	"platform/go-pkg/parser"
	"platform/go-pkg/response"
	"todo/types"

	"github.com/gin-gonic/gin"
)

func (h *TodoHandler) CreateTodoHandler(c *gin.Context) {
	log.Printf("hit service create todo with request: %v", c.Request)

	userID := c.GetString("user_id")

	var parsedBody types.DTOCreateTodo
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

	if parsedBody.Priority != "" {
		parsedBody.Priority = "medium"
	}

	log.Println("running controller")
	result, err := h.controller.CreateTodo(c.Request.Context(), &types.ReqCreateTodo{
		UserID: userID,
		Title: parsedBody.Title,
		Description: parsedBody.Description,
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

	log.Println("success create todo")
	response.Created(c, "success", result)
}