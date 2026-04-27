package repositories

import (
	"auth/types"
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type userRepository struct {
	db *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) IUserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(ctx context.Context, id uuid.UUID, email, passwordHash, fullName string) (err error) {
	query := `INSERT INTO users (id, email, password_hash, full_name) VALUES ($1, $2, $3, $4)`
	if _, err := r.db.Exec(ctx, query, id, email, passwordHash, fullName); err != nil {
		return err
	}
	return nil
}

func (r *userRepository) GetByEmail(ctx context.Context, email string) (user *types.User, err error) {
	query := `SELECT id, email, password_hash, full_name, avatar_url, created_at, updated_at FROM users WHERE email = $1 LIMIT 1`
	
	var u types.User
	if err := r.db.QueryRow(ctx, query, email).Scan(
		&u.ID,
		&u.Email,
		&u.PasswordHash,
		&u.FullName,
		&u.AvatarURL,
		&u.CreatedAt,
		&u.UpdatedAt,
	); err != nil {
		return nil, err
	}

	return &u, nil
}

func (r *userRepository) ExistsByID(ctx context.Context, id string) (exists bool, err error) {
	query := `SELECT EXISTS(SELECT 1 FROM users WHERE id = $1)`

	if err := r.db.QueryRow(ctx, query, id).Scan(&exists); err != nil {
		return exists, err
	}

	return exists, nil
}