package types

type (
	DTORegister struct {
		Email    string `json:"email" validate:"required,email,max=100"`
		Password string `json:"password" validate:"required,min=8"`
		FullName string `json:"full_name" validate:"required,min=2,max=100"`
	}
)

type (
	ReqRegister struct {
		Email    string
		Password string
		FullName string
	}

	ResRegister struct {
		User         User   `json:"user"`
		AccessToken  string `json:"access_token"`
		RefreshToken string `json:"refresh_token"`
		ExpiresIn    int    `json:"expires_in"`
	}
)