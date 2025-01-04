package dto

import "time"

type RegisterUserRequest struct {
	Name     string `validate:"required,min=3,max=100" json:"name"`
	Email    string `validate:"required,email" json:"email"`
	Password string `validate:"required,min=8" json:"password"`
	IsAdmin  bool   `json:"is_admin"`
}

type LoginUserRequest struct {
	Email    string `validate:"required,email" json:"email"`
	Password string `validate:"required" json:"password"`
}

type RenewAccessTokenRequest struct {
	RefreshToken string `validate:"required" json:"refresh_token"`
}

type RegisterUserResponse struct {
	Name    string `json:"name"`
	Email   string `json:"email"`
	IsAdmin bool   `json:"is_admin"`
}

type LoginUserResponse struct {
	AccessToken  string               `json:"access_token"`
	RefreshToken string               `json:"refresh_token"`
	SessionId    string               `json:"session_id"`
	User         RegisterUserResponse `json:"user"`
}

type RenewAccessTokenResponse struct {
	AccessToken          string    `json:"access_token"`
	AccessTokenExpiresAt time.Time `json:"access_toker_expires_at"`
}
