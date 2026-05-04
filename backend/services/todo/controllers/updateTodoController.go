package controllers

import (
	"context"
	"errors"
	"time"

	"todo/types"

	"github.com/google/uuid"
)

func (c *todoController) UpdateTodo(
	ctx context.Context,
	req *types.ReqUpdateTodo,
) (*types.ResUpdateTodo, error) {

	// --- VALIDATE UUID ---
	if _, err := uuid.Parse(req.TodoID); err != nil {
		return nil, errors.New("invalid todo id")
	}

	// --- CHECK OWNERSHIP ---
	ownerID, err := c.todoRepo.GetTodoOwner(ctx, req.TodoID)
	if err != nil {
		return nil, err
	}
	if ownerID == "" {
		return nil, errors.New("todo not found")
	}
	if ownerID != req.UserID {
		return nil, errors.New("forbidden")
	}

	// --- VALIDATE DUE DATE ---
	if req.DueDate != nil {
		if req.DueDate.Before(time.Now()) {
			return nil, errors.New("due_date must be in the future")
		}
	}

	// --- VALIDATE TAG LIMIT ---
	if len(req.TagIDs) > 10 {
		return nil, errors.New("max 10 tags allowed")
	}

	// --- CALL REPO ---
	res, err := c.todoRepo.UpdateTodo(ctx, req)
	if err != nil {
		return nil, err
	}

	return res, nil
}