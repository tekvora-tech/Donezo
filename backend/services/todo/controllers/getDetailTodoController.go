package controllers

import (
	"context"
	"errors"

	"todo/types"

	"github.com/google/uuid"
)

func (c *todoController) GetTodoDetail(
	ctx context.Context,
	req *types.ReqGetTodoDetail,
) (*types.ResGetTodoDetail, error) {

	// --- VALIDATE UUID ---
	if _, err := uuid.Parse(req.TodoID); err != nil {
		return nil, errors.New("invalid request: invalid todo id (must be UUID)")
	}

	// --- CALL REPO ---
	todo, err := c.todoRepo.GetTodoDetail(ctx, req)
	if err != nil {
		return nil, err
	}

	// not found → jangan bocorin apakah beda user atau tidak ada
	if todo == nil {
		return nil, errors.New("not found: todo not found")
	}

	return todo, nil
}