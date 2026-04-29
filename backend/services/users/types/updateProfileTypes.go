package types

import "time"

type (
	DTOUpdateProfile struct {
		FullName  string `json:"full_name" validate:"required,min=2,max=100"`
		AvatarURL string `json:"avatar_url" validate:"omitempty,url"`
	}
)

type (
	ReqUpdateProfile struct {
		UserID    string
		FullName  string
		AvatarURL string
	}

	ResUpdateProfile struct {
		ID        string     `json:"id"`
		Email     string     `json:"email"`
		FullName  string     `json:"full_name"`
		AvatarURL *string    `json:"avatar_url"`
		CreatedAt time.Time  `json:"created_at"`
		UpdatedAt *time.Time `json:"updated_at"`
	}
)