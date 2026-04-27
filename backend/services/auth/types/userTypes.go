package types

import "time"

type (
	User struct {
		ID           string `json:"id"`
		Email        string `json:"email"`
		PasswordHash string `json:"password_hash,omitempty"`
		FullName     string `json:"full_name"`
		AvatarURL    *string `json:"avatar_url"`
		CreatedAt    *time.Time `json:"created_at,omitempty"`
		UpdatedAt *time.Time `json:"updated_at,omitempty"`
	}
)