package service

import (
	"github.com/sourav014/go-auth-service/dto"
)

type AuthService interface {
	RegisterUser(registerUserRequest dto.RegisterUserRequest) (dto.RegisterUserResponse, error)
	LoginUser(loginUserRequest dto.LoginUserRequest) (dto.LoginUserResponse, error)
	RenewAccessToken(renewAccessTokenRequest dto.RenewAccessTokenRequest) (dto.RenewAccessTokenResponse, error)
	RevokeToken(sessionId string) error
}
