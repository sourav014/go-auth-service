package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sourav014/go-auth-service/dto"
	"github.com/sourav014/go-auth-service/models"
	"github.com/sourav014/go-auth-service/service"
)

type AuthController struct {
	authService service.AuthService
}

func NewAuthController(authService service.AuthService) *AuthController {
	return &AuthController{authService: authService}
}

func (authController *AuthController) RegisterUser(ctx *gin.Context) {
	var registerUserRequest dto.RegisterUserRequest
	if err := ctx.ShouldBindJSON(&registerUserRequest); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	registerUserResponse, err := authController.authService.RegisterUser(registerUserRequest)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusCreated, registerUserResponse)
}

func (authController *AuthController) LoginUser(ctx *gin.Context) {
	var loginUserRequest dto.LoginUserRequest
	if err := ctx.ShouldBindJSON(&loginUserRequest); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	loginUserResponse, err := authController.authService.LoginUser(loginUserRequest)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, loginUserResponse)
}

func (authController *AuthController) RenewAccessToken(ctx *gin.Context) {
	var renewAccessTokenRequest dto.RenewAccessTokenRequest
	if err := ctx.ShouldBindJSON(&renewAccessTokenRequest); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	renewAccessTokenResponse, err := authController.authService.RenewAccessToken(renewAccessTokenRequest)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, renewAccessTokenResponse)

}

func (authController *AuthController) RevokeToken(ctx *gin.Context) {
	sessionId := ctx.Param("id")

	err := authController.authService.RevokeToken(sessionId)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "token revoked successfully",
	})
}

func (authController *AuthController) GetUserProfile(ctx *gin.Context) {
	currentUser, ok := ctx.MustGet("currentUser").(models.User)
	if !ok {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "internal server error",
		})
		return
	}

	userProfile := dto.RegisterUserResponse{
		Name:    currentUser.Name,
		Email:   currentUser.Email,
		IsAdmin: currentUser.IsAdmin,
	}

	ctx.JSON(http.StatusOK, userProfile)
}
