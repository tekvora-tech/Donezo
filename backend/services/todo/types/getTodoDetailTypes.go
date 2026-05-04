package types

import "time"

type (
	ReqGetTodoDetail struct {
		UserID string
		TodoID string
	}

	ResGetTodoDetail struct {
		ID                    string     `json:"id"`
		UserID string `json:"user_id"`
		Title                 string     `json:"title"`
		Description           string     `json:"description"`
		Status                string     `json:"status"`
		Priority              string     `json:"priority"`
		DueDate               time.Time  `json:"due_date"`
		Tags                  []Tag      `json:"tags"`
		SubTasks []SubTask `json:"sub_tasks"`
		CreatedAt             time.Time  `json:"created_at"`
		UpdatedAt             *time.Time `json:"updated_at"`
	}
)