package repositories

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"
	"todo/types"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type todoRepository struct {
	db *pgxpool.Pool
}

func NewTodoRepository(db *pgxpool.Pool) ITodoRepository {
	return &todoRepository{db: db}
}

func (r *todoRepository) CreateTodo(
	ctx context.Context,
	todoID, userID, title, description, status, priority string,
	dueDate time.Time,
) (*types.SimpleTodo, error) {

	query := `
		INSERT INTO todos (
			id, user_id, title, description, status, priority, due_date
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7
		)
		RETURNING 
			id, user_id, title, description, status, priority, due_date, created_at, updated_at
	`

	var res types.SimpleTodo

	err := r.db.QueryRow(
		ctx,
		query,
		todoID,
		userID,
		title,
		description,
		status,
		priority,
		dueDate,
	).Scan(
		&res.ID,
		&res.UserID,
		&res.Title,
		&res.Description,
		&res.Status,
		&res.Priority,
		&res.DueDate,
		&res.CreatedAt,
		&res.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &res, nil
}

func (r *todoRepository) ListTodos(
	ctx context.Context,
	req *types.ReqListTodo,
) ([]types.Todo, int, error) {

	var (
		args       []interface{}
		conditions []string
		argPos     = 1
	)

	// BASE
	baseQuery := `
	FROM todos t
	LEFT JOIN todo_tags tt ON t.id = tt.todo_id
	LEFT JOIN tags tg ON tt.tag_id = tg.id
	LEFT JOIN sub_tasks st ON t.id = st.todo_id
	`

	// WHERE
	conditions = append(conditions, fmt.Sprintf("t.user_id = $%d", argPos))
	args = append(args, req.UserID)
	argPos++

	// FILTERS
	if req.Status != "" {
		conditions = append(conditions, fmt.Sprintf("t.status = $%d", argPos))
		args = append(args, req.Status)
		argPos++
	}

	if req.Priority != "" {
		conditions = append(conditions, fmt.Sprintf("t.priority = $%d", argPos))
		args = append(args, req.Priority)
		argPos++
	}

	if req.TagID != "" {
		conditions = append(conditions, fmt.Sprintf("tt.tag_id = $%d", argPos))
		args = append(args, req.TagID)
		argPos++
	}

	if req.Search != "" {
		conditions = append(conditions, fmt.Sprintf("t.title ILIKE $%d", argPos))
		args = append(args, "%"+req.Search+"%")
		argPos++
	}

	if req.DueDateFrom != nil {
		conditions = append(conditions, fmt.Sprintf("t.due_date >= $%d", argPos))
		args = append(args, *req.DueDateFrom)
		argPos++
	}

	if req.DueDateTo != nil {
		conditions = append(conditions, fmt.Sprintf("t.due_date <= $%d", argPos))
		args = append(args, *req.DueDateTo)
		argPos++
	}

	whereClause := "WHERE " + strings.Join(conditions, " AND ")

	// ======================
	// COUNT QUERY
	// ======================
	countQuery := "SELECT COUNT(DISTINCT t.id) " + baseQuery + " " + whereClause

	var totalRecords int
	err := r.db.QueryRow(ctx, countQuery, args...).Scan(&totalRecords)
	if err != nil {
		return nil, 0, err
	}

	// ======================
	// MAIN QUERY
	// ======================
	offset := (req.Page - 1) * req.PageSize

	query := fmt.Sprintf(`
	SELECT 
		t.id,
		t.title,
		t.description,
		t.status,
		t.priority,
		t.due_date,
		t.created_at,
		t.updated_at,

		-- TAGS
		COALESCE(
			json_agg(DISTINCT jsonb_build_object(
				'id', tg.id,
				'name', tg.name,
				'color', tg.color
			)) FILTER (WHERE tg.id IS NOT NULL),
			'[]'
		) as tags,

		-- SUB TASK COUNT
		COUNT(DISTINCT st.id) as sub_task_count,
		COUNT(DISTINCT CASE WHEN st.is_completed THEN st.id END) as completed_sub_task_count

	%s
	%s
	GROUP BY t.id
	ORDER BY t.%s %s
	LIMIT %d OFFSET %d
	`, baseQuery, whereClause, req.SortBy, strings.ToUpper(req.SortOrder), req.PageSize, offset)

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var todos []types.Todo

	for rows.Next() {
		var t types.Todo
		var tagsBytes []byte

		err := rows.Scan(
			&t.ID,
			&t.Title,
			&t.Description,
			&t.Status,
			&t.Priority,
			&t.DueDate,
			&t.CreatedAt,
			&t.UpdatedAt,
			&tagsBytes,
			&t.SubTaskCount,
			&t.CompletedSubTaskCount,
		)
		if err != nil {
			return nil, 0, err
		}

		// decode JSON tags
		if err := json.Unmarshal(tagsBytes, &t.Tags); err != nil {
			return nil, 0, err
		}

		todos = append(todos, t)
	}

	return todos, totalRecords, nil
}

func (r *todoRepository) GetTodoDetail(
	ctx context.Context,
	req *types.ReqGetTodoDetail,
) (*types.ResGetTodoDetail, error) {

	query := `
	SELECT 
		t.id,
		t.user_id,
		t.title,
		t.description,
		t.status,
		t.priority,
		t.due_date,
		t.created_at,
		t.updated_at,

		COALESCE(
			json_agg(DISTINCT jsonb_build_object(
				'id', tg.id,
				'name', tg.name,
				'color', tg.color
			)) FILTER (WHERE tg.id IS NOT NULL),
			'[]'
		) as tags,

		COALESCE(
			json_agg(DISTINCT jsonb_build_object(
				'id', st.id,
				'title', st.title,
				'is_completed', st.is_completed,
				'created_at', st.created_at
			)) FILTER (WHERE st.id IS NOT NULL),
			'[]'
		) as sub_tasks

	FROM todos t
	LEFT JOIN todo_tags tt ON t.id = tt.todo_id
	LEFT JOIN tags tg ON tt.tag_id = tg.id
	LEFT JOIN sub_tasks st ON t.id = st.todo_id
	WHERE t.id = $1 AND t.user_id = $2
	GROUP BY t.id
	`

	var (
		res       types.ResGetTodoDetail
		tagsJSON  []byte
		subJSON   []byte
	)

	err := r.db.QueryRow(ctx, query, req.TodoID, req.UserID).Scan(
		&res.ID,
		&res.UserID,
		&res.Title,
		&res.Description,
		&res.Status,
		&res.Priority,
		&res.DueDate,
		&res.CreatedAt,
		&res.UpdatedAt,
		&tagsJSON,
		&subJSON,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil // biar controller yg handle 404
		}
		return nil, err
	}

	// decode tags
	if err := json.Unmarshal(tagsJSON, &res.Tags); err != nil {
		return nil, err
	}

	// decode sub tasks
	if err := json.Unmarshal(subJSON, &res.SubTasks); err != nil {
		return nil, err
	}

	return &res, nil
}

// ========================
// GET OWNER
// ========================
func (r *todoRepository) GetTodoOwner(ctx context.Context, todoID string) (string, error) {
	var userID string

	err := r.db.QueryRow(ctx,
		`SELECT user_id FROM todos WHERE id = $1`,
		todoID,
	).Scan(&userID)

	if err != nil {
		if err == pgx.ErrNoRows {
			return "", nil
		}
		return "", err
	}

	return userID, nil
}

// ========================
// UPDATE TODO
// ========================
func (r *todoRepository) UpdateTodo(
	ctx context.Context,
	req *types.ReqUpdateTodo,
) (*types.ResUpdateTodo, error) {

	tx, err := r.db.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	// ========================
	// DYNAMIC UPDATE
	// ========================
	var (
		setClauses []string
		args       []interface{}
		argPos     = 1
	)

	if req.Title != "" {
		setClauses = append(setClauses, fmt.Sprintf("title = $%d", argPos))
		args = append(args, req.Title)
		argPos++
	}

	if req.Description != "" {
		setClauses = append(setClauses, fmt.Sprintf("description = $%d", argPos))
		args = append(args, req.Description)
		argPos++
	}

	if req.Status != "" {
		setClauses = append(setClauses, fmt.Sprintf("status = $%d", argPos))
		args = append(args, req.Status)
		argPos++
	}

	if req.Priority != "" {
		setClauses = append(setClauses, fmt.Sprintf("priority = $%d", argPos))
		args = append(args, req.Priority)
		argPos++
	}

	if req.DueDate != nil {
		setClauses = append(setClauses, fmt.Sprintf("due_date = $%d", argPos))
		args = append(args, *req.DueDate)
		argPos++
	}

	if len(setClauses) > 0 {
		setClauses = append(setClauses, "updated_at = now()")

		query := fmt.Sprintf(`
			UPDATE todos
			SET %s
			WHERE id = $%d
		`, strings.Join(setClauses, ", "), argPos)

		args = append(args, req.TodoID)

		_, err := tx.Exec(ctx, query, args...)
		if err != nil {
			return nil, err
		}
	}

	// ========================
	// TAG REPLACE
	// ========================
	if req.TagIDs != nil {

		// delete old
		_, err := tx.Exec(ctx,
			`DELETE FROM todo_tags WHERE todo_id = $1`,
			req.TodoID,
		)
		if err != nil {
			return nil, err
		}

		// insert new
		for _, tagID := range req.TagIDs {

			// OPTIONAL: validate tag ownership
			var exists bool
			err := tx.QueryRow(ctx,
				`SELECT EXISTS (
					SELECT 1 FROM tags WHERE id = $1 AND user_id = $2
				)`,
				tagID, req.UserID,
			).Scan(&exists)

			if err != nil {
				return nil, err
			}
			if !exists {
				return nil, fmt.Errorf("tag %s not found", tagID)
			}

			_, err = tx.Exec(ctx,
				`INSERT INTO todo_tags (todo_id, tag_id) VALUES ($1, $2)`,
				req.TodoID, tagID,
			)
			if err != nil {
				return nil, err
			}
		}
	}

	// ========================
	// GET UPDATED DATA
	// ========================
	query := `
	SELECT 
		t.id,
		t.user_id,
		t.title,
		t.description,
		t.status,
		t.priority,
		t.due_date,
		t.created_at,
		t.updated_at,

		COALESCE(
			json_agg(DISTINCT jsonb_build_object(
				'id', tg.id,
				'name', tg.name,
				'color', tg.color
			)) FILTER (WHERE tg.id IS NOT NULL),
			'[]'
		) as tags,

		COALESCE(
			json_agg(DISTINCT jsonb_build_object(
				'id', st.id,
				'title', st.title,
				'is_completed', st.is_completed,
				'created_at', st.created_at
			)) FILTER (WHERE st.id IS NOT NULL),
			'[]'
		) as sub_tasks

	FROM todos t
	LEFT JOIN todo_tags tt ON t.id = tt.todo_id
	LEFT JOIN tags tg ON tt.tag_id = tg.id
	LEFT JOIN sub_tasks st ON t.id = st.todo_id
	WHERE t.id = $1
	GROUP BY t.id
	`

	var (
		res     types.ResUpdateTodo
		tags    []byte
		subtask []byte
	)

	err = tx.QueryRow(ctx, query, req.TodoID).Scan(
		&res.ID,
		&res.UserID,
		&res.Title,
		&res.Description,
		&res.Status,
		&res.Priority,
		&res.DueDate,
		&res.CreatedAt,
		&res.UpdatedAt,
		&tags,
		&subtask,
	)
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(tags, &res.Tags); err != nil {
		return nil, err
	}
	if err := json.Unmarshal(subtask, &res.SubTasks); err != nil {
		return nil, err
	}

	// COMMIT
	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}

	return &res, nil
}

func (r *todoRepository) DeleteTodo(
	ctx context.Context,
	todoID string,
	userID string,
) error {

	cmdTag, err := r.db.Exec(ctx,
		`DELETE FROM todos WHERE id = $1 AND user_id = $2`,
		todoID,
		userID,
	)
	if err != nil {
		return err
	}

	// safety check (should not happen karena sudah check ownership di controller)
	if cmdTag.RowsAffected() == 0 {
		return errors.New("todo not found")
	}

	return nil
}