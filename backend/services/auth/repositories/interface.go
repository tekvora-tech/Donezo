package repositories

import (
	"auth/types"
	"context"

	"github.com/google/uuid"
)

type (
	IUserRepository interface {
		Create(ctx context.Context, id uuid.UUID, email, passwordHash, fullName string) (err error)
		GetByEmail(ctx context.Context, email string) (user *types.User, err error)
		ExistsByID(ctx context.Context, id string) (exists bool, err error)
	}
)