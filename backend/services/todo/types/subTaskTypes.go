package types

import "time"

type (
	SubTask struct {
		ID          string `json:"id"`
		Title       string `json:"title"`
		IsCompleted bool   `json:"is_completed"`
		CreatedAt   time.Time `json:"created_at"`
	}
)