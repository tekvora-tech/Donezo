package types

type (
	ReqRefreshToken struct {
		RefreshToken string
	}

	ResRefreshToken struct {
		AccessToken string `json:"access_token"`
		ExpiresIn   int    `json:"expires_in"`
	}
)