package controllers

import (
	"context"
	"fmt"
	"todo/types"

	"github.com/google/uuid"
)

func (c *todoController) CreateTodo(ctx context.Context, req *types.ReqCreateTodo) (res *types.ResCreateTodo, err error) {
	tags, allValid, err := c.tagRepo.GetTagsAndValidate(ctx, req.TagIDs)
	if err != nil {
		return nil, err	
	}

	if !allValid {
		return nil, fmt.Errorf("invalid request: invalid tags ids")
	}

	todoID := uuid.NewString()
	status := defaultStatus

	todo, err := c.todoRepo.CreateTodo(ctx, todoID, req.UserID, req.Title, req.Description, status, req.Priority, *req.DueDate)
	if err != nil {
		return nil, err
	}

	rows := [][]interface{}{}

	for _, tagID := range req.TagIDs {
		rows = append(rows, []interface{}{todoID, tagID})
	}

	if err := c.tagRepo.AssignTags(ctx, rows); err != nil {
		return nil, err
	}

	return &types.ResCreateTodo{
		ID: todo.ID,
		UserID: todo.UserID,
		Title: todo.Title,
		Description: todo.Description,
		Status: todo.Status,
		Priority: todo.Priority,
		DueDate: todo.DueDate,
		Tags: tags,
		SubTaks: []string{},
		CreatedAt: todo.CreatedAt,
		UpdatedAt: todo.UpdatedAt,
	}, nil
}