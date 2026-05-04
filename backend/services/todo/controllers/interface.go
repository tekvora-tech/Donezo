package controllers

import (
	"context"
	"todo/types"
)

var (
	defaultStatus = "pending"
)

type (
	ITodoController interface {
		CreateTodo(ctx context.Context, req *types.ReqCreateTodo) (res *types.ResCreateTodo, err error)
		ListTodo(ctx context.Context, req *types.ReqListTodo) (res *types.ResListTodo, err error)
		GetTodoDetail(ctx context.Context, req *types.ReqGetTodoDetail) (res *types.ResGetTodoDetail, err error)
		UpdateTodo(ctx context.Context, req *types.ReqUpdateTodo) (res *types.ResUpdateTodo, err error)
		DeleteTodo(ctx context.Context, req *types.ReqDeleteTodo) (res *types.ResDeleteTodo, err error)
	}
)