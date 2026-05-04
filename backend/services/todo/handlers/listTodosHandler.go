package handlers

import (
	"log"
	"platform/go-pkg/parser"
	"platform/go-pkg/response"
	"strconv"
	"time"
	"todo/types"

	"github.com/gin-gonic/gin"
)

func (h *TodoHandler) ListTodosHandler(c *gin.Context) {
	log.Printf("hit service list todos with request: %v", c.Request)

	userID := c.GetString("user_id")

	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil || page <= 0 {
		log.Printf("invalid page: %v", err)
		response.BadRequest(c, "invalid page", err.Error())
		return
	}

	pageSize, err := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	if err != nil || (pageSize <= 0  || pageSize > 100) {
		log.Printf("invalid page size: %v", err)
		response.BadRequest(c, "invalid page size", err.Error())
		return
	}

	var dueDateFrom *time.Time
	queryDueDateFrom := c.Query("due_date_from")
	if queryDueDateFrom != "" {
		parsed, err := time.Parse(time.RFC3339, queryDueDateFrom)
		if err != nil {
			log.Printf("failed parse due date from: %v", err)
			response.BadRequest(c, "failed parse due date from", err.Error())
			return
		}
		dueDateFrom = &parsed
	} else {
		dueDateFrom = nil
	}
	
	var dueDateTo *time.Time
	queryDueDateTo := c.Query("due_date_to")
	if queryDueDateTo != "" {
		parsed, err := time.Parse(time.RFC3339, queryDueDateTo)
		if err != nil {
			log.Printf("failed parse due date to: %v", err)
			response.BadRequest(c, "failed parse due date to", err.Error())
			return
		}
		dueDateTo = &parsed
	} else {
		dueDateTo = nil
	}

	log.Println("running controller")
	result, err := h.controller.ListTodo(c.Request.Context(), &types.ReqListTodo{
		UserID: userID,
		Page: page,
		PageSize: pageSize,
		Status: c.Query("status"),
		Priority: c.Query("priority"),
		TagID: c.Query("tag_id"),
		Search: c.Query("search"),
		DueDateFrom: dueDateFrom,
		DueDateTo: dueDateTo,
		SortBy: c.DefaultQuery("sort_by", "created_at"),
		SortOrder: c.DefaultQuery("sort_order", "desc"),
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

	log.Println("success list todos")
	response.OK(c, "success", result)
}