package repositories

import (
	"context"
	"fmt"
	"users/types"

	"github.com/jackc/pgx/v5/pgxpool"
)

type userRepository struct {
	db *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) IUserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Update(ctx context.Context, userID string, fullName, avatarURL string) (err error) {
	query := `UPDATE users SET full_name = $1, avatar_url = $2 WHERE id = $3`
	res, err := r.db.Exec(ctx, query, fullName, avatarURL, userID)
	if err != nil {
		return err
	}

	if res.RowsAffected() == 0 {
		return fmt.Errorf("not found: data tidak ditemukan")
	}

	return nil
}

func (r *userRepository) GetByID(ctx context.Context, userID string) (user *types.User, err error) {
	query := `SELECT id, email, full_name, avatar_url, created_at, updated_at FROM users WHERE id = $1 LIMIT 1`
	
	var u types.User
	if err := r.db.QueryRow(ctx, query, userID).Scan(
		&u.ID,
		&u.Email,
		&u.FullName,
		&u.AvatarURL,
		&u.CreatedAt,
		&u.UpdatedAt,
	); err != nil {
		return nil, err
	}

	return &u, nil
}