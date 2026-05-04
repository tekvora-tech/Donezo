package types

import "time"

type (
	ReqListTodo struct {
		UserID   string
		Page     int
		PageSize int
		Status   string
		Priority string
		TagID    string
		Search   string
		DueDateFrom  *time.Time
		DueDateTo *time.Time
		SortBy string
		SortOrder string
	}

	ResListTodo struct {
		Data []Todo `json:"data"`
		Metadata Metadata `json:"metadata"`
	}
)