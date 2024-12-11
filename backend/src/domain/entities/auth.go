package entities

type TokenResponseGithub struct {
	AccessToken string `json:"access_token"`
	Scope       string `json:"scope"`
	TokenType   string `json:"token_type"`
	Error       string `json:"error"`
	ErrorDesc   string `json:"error_description"`
	JWT         string `json:"jwt"`
}

type UserGithub struct {
	Username string `json:"login"`
	UserID   string `json:"node_id"`
	ImageURL string `json:"avatar_url"`
	Email    string `json:"email"`
}
