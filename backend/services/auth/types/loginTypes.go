package types

type (
	DTOLogin struct {
		Email    string `json:"email" validate:"required,email,max=100"`
		Password string `json:"password" validate:"required,min=8"`
	}
)

type (
	ReqLogin struct {
		Email    string
		Password string
	}

	ResLogin struct {
		User         User   `json:"user"`
		AccessToken  string `json:"access_token"`
		RefreshToken string `json:"refresh_token"`
		ExpiresIn    int    `json:"expires_in"`
	}
)