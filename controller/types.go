package controller

type RegisterUserInput struct {
	Name     string `json:"name" binding:"required,min=3,max=100"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
	IsAdmin  bool   `json:"is_admin" binding:"required"`
}

type LoginUserInput struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type RenewAccessTokenInput struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

type RegisterUserResp struct {
	Name    string `json:"name"`
	Email   string `json:"email"`
	IsAdmin bool   `json:"is_admin"`
}

type LoginUserResp struct {
	AccessToken  string           `json:"access_token"`
	RefreshToken string           `json:"refresh_token"`
	SessionId    string           `json:"session_id"`
	User         RegisterUserResp `json:"user"`
}
