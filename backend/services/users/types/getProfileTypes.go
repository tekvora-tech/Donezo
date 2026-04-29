package types

import "time"

type (
	ReqGetProfile struct {
		UserID string
	}

	ResGetProfile struct {
		ID           string     `json:"id"`
		Email        string     `json:"email"`
		FullName     string     `json:"full_name"`
		AvatarURL    *string    `json:"avatar_url"`
		CreatedAt    time.Time `json:"created_at"`
		UpdatedAt    *time.Time `json:"updated_at"`
	}
)