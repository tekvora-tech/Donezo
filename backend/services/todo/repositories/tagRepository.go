package repositories

import (
	"context"
	"todo/types"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type tagRepository struct {
	db *pgxpool.Pool
}

func NewTagRepository(db *pgxpool.Pool) ITagRepository {
	return &tagRepository{db: db}
}

func (r *tagRepository) GetTagsAndValidate(ctx context.Context, ids []string) ([]types.Tag, bool, error) {
	query := `SELECT id, name, color FROM tags WHERE id = ANY($1)`

	rows, err := r.db.Query(ctx, query, ids)
	if err != nil {
		return nil, false, err
	}
	defer rows.Close()

	var tags []types.Tag
	for rows.Next() {
		var t types.Tag
		if err := rows.Scan(&t.ID, &t.Name, &t.Color); err != nil {
			return nil, false, err
		}
		tags = append(tags, t)
	}

	if err := rows.Err(); err != nil {
		return nil, false, err
	}

	// valid kalau jumlah yang ketemu == jumlah input
	allValid := len(tags) == len(ids)

	return tags, allValid, nil
}

func (r *tagRepository) AssignTags(ctx context.Context, data [][]interface{}) (err error) {
	if _, err := r.db.CopyFrom(
		ctx,
		pgx.Identifier{"todo_tags"},
		[]string{"todo_id", "tag_id"},
		pgx.CopyFromRows(data),
	); err != nil {
		return err
	}

	return nil
}