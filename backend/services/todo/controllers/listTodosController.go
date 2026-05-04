package controllers

import (
	"context"
	"errors"
	"math"
	"strings"
	"todo/types"
)

var validSortFields = map[string]bool{
	"created_at": true,
	"updated_at": true,
	"due_date":   true,
	"priority":   true,
	"title":      true,
}

var validSortOrder = map[string]bool{
	"asc":  true,
	"desc": true,
}

func (c *todoController) ListTodo(ctx context.Context, req *types.ReqListTodo) (*types.ResListTodo, error) {
	// --- VALIDATION ---
	if req.Page < 1 {
		req.Page = 1
	}

	if req.PageSize < 1 {
		req.PageSize = 10
	}
	if req.PageSize > 100 {
		req.PageSize = 100
	}

	// validate sort
	if !validSortFields[req.SortBy] {
		req.SortBy = "created_at"
	}
	if !validSortOrder[strings.ToLower(req.SortOrder)] {
		req.SortOrder = "desc"
	}

	// sanitize search
	if len(req.Search) > 100 {
		return nil, errors.New("invalid request: search too long (max 100 chars)")
	}

	// --- CALL REPOSITORY ---
	todos, totalRecords, err := c.todoRepo.ListTodos(ctx, req)
	if err != nil {
		return nil, err
	}

	// --- METADATA ---
	totalPages := int(math.Ceil(float64(totalRecords) / float64(req.PageSize)))

	res := &types.ResListTodo{
		Data: todos,
		Metadata: types.Metadata{
			CurrentPage:  req.Page,
			PageSize:     req.PageSize,
			TotalPages:   totalPages,
			TotalRecords: totalRecords,
			HasNext:      req.Page < totalPages,
			HasPrev:      req.Page > 1,
		},
	}

	return res, nil
}