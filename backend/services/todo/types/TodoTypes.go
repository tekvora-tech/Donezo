package types

import "time"

type (
	Todo struct {
		ID          string `json:"id"`
		Title       string `json:"title"`
		Description string `json:"description"`
		Status      string `json:"status"`
		Priority    string `json:"priority"`
		DueDate     time.Time `json:"due_date"`
		Tags []Tag `json:"tags"`
		SubTaskCount int `json:"sub_task_count"`
		CompletedSubTaskCount int `json:"completed_sub_task_count"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt *time.Time `json:"updated_at"`
	}

	SimpleTodo struct {
		ID          string    `json:"id"`
		UserID      string    `json:"user_id"`
		Title       string    `json:"title"`
		Description string    `json:"description"`
		Status      string    `json:"status"`
		Priority    string    `json:"priority"`
		DueDate     time.Time `json:"due_date"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt *time.Time `json:"updated_at"`
	}
)