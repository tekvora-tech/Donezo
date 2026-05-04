package repositories

import (
	"context"
	"time"
	"todo/types"
)

type (
	ITodoRepository interface {
		CreateTodo(ctx context.Context, todoID, userID, title, description, status, priority string, dueDate time.Time) (*types.SimpleTodo, error)
		ListTodos(ctx context.Context, req *types.ReqListTodo) ([]types.Todo, int, error)
		GetTodoDetail(ctx context.Context, req *types.ReqGetTodoDetail) (*types.ResGetTodoDetail, error)
		GetTodoOwner(ctx context.Context, todoID string) (string, error)
		UpdateTodo(ctx context.Context, req *types.ReqUpdateTodo) (*types.ResUpdateTodo, error)
		DeleteTodo(ctx context.Context, todoID string, userID string) error
	}

	ITagRepository interface {
		GetTagsAndValidate(ctx context.Context, ids []string) ([]types.Tag, bool, error)
		AssignTags(ctx context.Context, data [][]interface{}) (err error)
	}
)