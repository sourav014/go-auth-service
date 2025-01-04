package service

import (
	"errors"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/sourav014/go-auth-service/dto"
	"github.com/sourav014/go-auth-service/models"
	"github.com/sourav014/go-auth-service/repository"
	JwtToken "github.com/sourav014/go-auth-service/token"
	"github.com/sourav014/go-auth-service/util"
)

type AuthServiceImpl struct {
	SessionsRepository repository.SessionsRepository
	UsersReposity      repository.UsersRepository
	JWTMaker           *JwtToken.JWTMaker
	Validate           *validator.Validate
}

func NewAuthServiceImpl(sessionRepository repository.SessionsRepository, userRepository repository.UsersRepository, jwtMaker *JwtToken.JWTMaker, validate *validator.Validate) AuthService {
	return &AuthServiceImpl{
		SessionsRepository: sessionRepository,
		UsersReposity:      userRepository,
		JWTMaker:           jwtMaker,
		Validate:           validate,
	}
}

func (auth *AuthServiceImpl) RegisterUser(registerUserRequest dto.RegisterUserRequest) (dto.RegisterUserResponse, error) {
	var registerUserResponse dto.RegisterUserResponse
	err := auth.Validate.Struct(registerUserRequest)
	if err != nil {
		return registerUserResponse, err
	}

	user, err := auth.UsersReposity.FindByEmail(registerUserRequest.Email)
	if err != nil {
		return registerUserResponse, err
	}

	if user.ID != 0 {
		return registerUserResponse, errors.New("email already exists")
	}

	hashedPassword, err := util.GenerateHashString(registerUserRequest.Password)
	if err != nil {
		return registerUserResponse, err
	}

	newUser := models.User{
		Name:     registerUserRequest.Name,
		Email:    registerUserRequest.Email,
		Password: hashedPassword,
		IsAdmin:  registerUserRequest.IsAdmin,
	}
	auth.UsersReposity.Create(newUser)

	registerUserResponse = dto.RegisterUserResponse{
		Name:    newUser.Name,
		Email:   newUser.Email,
		IsAdmin: newUser.IsAdmin,
	}

	return registerUserResponse, nil
}

func (auth *AuthServiceImpl) LoginUser(loginUserRequest dto.LoginUserRequest) (dto.LoginUserResponse, error) {
	var loginUserResponse dto.LoginUserResponse
	err := auth.Validate.Struct(loginUserRequest)
	if err != nil {
		return loginUserResponse, err
	}

	user, err := auth.UsersReposity.FindByEmail(loginUserRequest.Email)
	if err != nil {
		return loginUserResponse, err
	}

	err = util.CompareHashString(user.Password, loginUserRequest.Password)
	if err != nil {
		return loginUserResponse, errors.New("invalid credentials")
	}

	accessToken, _, err := auth.JWTMaker.CreateToken(user.ID, user.IsAdmin, time.Minute*15)
	if err != nil {
		return loginUserResponse, errors.New("failed to generate token")
	}

	refreshToken, refreshClaims, err := auth.JWTMaker.CreateToken(user.ID, user.IsAdmin, time.Hour*24)
	if err != nil {
		return loginUserResponse, errors.New("failed to generate token")
	}
	session := models.Session{
		ID:           refreshClaims.RegisteredClaims.ID,
		UserEmail:    user.Email,
		RefreshToken: refreshToken,
		IsRevoked:    false,
		ExpiresAt:    refreshClaims.RegisteredClaims.ExpiresAt.Time,
	}

	auth.SessionsRepository.Create(session)

	loginUserResponse = dto.LoginUserResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		SessionId:    session.ID,
		User: dto.RegisterUserResponse{
			Name:    user.Name,
			Email:   user.Email,
			IsAdmin: user.IsAdmin,
		},
	}

	return loginUserResponse, nil
}

func (auth *AuthServiceImpl) RenewAccessToken(renewAccessTokenRequest dto.RenewAccessTokenRequest) (dto.RenewAccessTokenResponse, error) {
	var renewAccessTokenResponse dto.RenewAccessTokenResponse
	err := auth.Validate.Struct(renewAccessTokenRequest)
	if err != nil {
		return renewAccessTokenResponse, err
	}
	refreshClaims, err := auth.JWTMaker.VerifyToken(renewAccessTokenRequest.RefreshToken)
	if err != nil {
		return renewAccessTokenResponse, err
	}

	session, err := auth.SessionsRepository.FindById(refreshClaims.RegisteredClaims.ID)
	if err != nil {
		return renewAccessTokenResponse, err
	}

	if session.IsRevoked {
		return renewAccessTokenResponse, errors.New("session has revoked")
	}

	accessToken, accessClaims, err := auth.JWTMaker.CreateToken(refreshClaims.ID, refreshClaims.IsAdmin, time.Minute*15)
	if err != nil {
		return renewAccessTokenResponse, errors.New("internal server error")
	}

	renewAccessTokenResponse = dto.RenewAccessTokenResponse{
		AccessToken:          accessToken,
		AccessTokenExpiresAt: accessClaims.RegisteredClaims.ExpiresAt.Time,
	}

	return renewAccessTokenResponse, nil
}

func (auth *AuthServiceImpl) RevokeToken(sessionId string) error {
	_, err := auth.SessionsRepository.FindById(sessionId)
	if err != nil {
		return err
	}

	auth.SessionsRepository.Update(sessionId)
	return nil
}
