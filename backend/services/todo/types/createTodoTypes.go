package types

import "time"

type (
	DTOCreateTodo struct {
		Title       string `json:"title" validate:"required,min=3,max=255"`
		Description string `json:"description" validate:"omitempty,max=5000"`
		Priority    string `json:"priority" validate:"omitempty,oneof=low medium high urgent"`
		DueDate     *time.Time `json:"due_date" validate:"omitempty,gt"`
		TagIDs 	[]string `jso:"tag_ids" validator:"omitempty,dive,uuid"`
	}
)

type (
	ReqCreateTodo struct {
		UserID string
		Title string
		Description string
		Priority string
		DueDate *time.Time
		TagIDs []string
	}

	ResCreateTodo struct {
		ID string `json:"id"`
		UserID string `json:"user_id"`
		Title string `json:"title"`
		Description string `json:"description"`
		Status string `json:"status"`
		Priority string `json:"priority"`
		DueDate time.Time `json:"due_date"`
		Tags []Tag `json:"tags"`
		SubTaks []string
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt *time.Time `json:"updated_at"`
	}
)