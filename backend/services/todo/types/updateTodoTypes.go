package types

import "time"

type (
	DTOUpdateTodo struct {
		Title       string    `json:"title" validate:"omitempty,min=3,max=255"`
		Description string    `json:"description" validate:"omitempty,max=5000"`
		Status string `json:"status" validate:"omitempty,oneof=pending in_progress completed cancelled"`
		Priority    string    `json:"priority" validate:"omitempty,oneof=low medium high urgent"`
		DueDate     *time.Time `json:"due_date" validate:"omitempty,gt"`
		TagIDs []string `json:"tag_ids" validate:"omitempty,dive,uuid"`
	}
)

type (
	ReqUpdateTodo struct {
		UserID string
		TodoID string
		Title string
		Description string
		Status string
		Priority string
		DueDate *time.Time
		TagIDs []string
	}

	ResUpdateTodo struct {
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