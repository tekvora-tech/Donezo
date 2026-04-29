package repositories

import (
	"context"
	"users/types"
)

type (
	IUserRepository interface {
		GetByID(ctx context.Context, userID string) (user *types.User, err error)
		Update(ctx context.Context, userID string, fullName, avatarURL string) (err error)
	}
)