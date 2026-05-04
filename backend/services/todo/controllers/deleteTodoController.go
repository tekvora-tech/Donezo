package controllers

import (
	"context"
	"errors"

	"todo/types"

	"github.com/google/uuid"
)

func (c *todoController) DeleteTodo(
	ctx context.Context,
	req *types.ReqDeleteTodo,
) (*types.ResDeleteTodo, error) {

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

	// --- DELETE ---
	err = c.todoRepo.DeleteTodo(ctx, req.TodoID, req.UserID)
	if err != nil {
		return nil, err
	}

	return &types.ResDeleteTodo{}, nil
}