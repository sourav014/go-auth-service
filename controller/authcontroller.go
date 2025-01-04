package controller

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sourav014/go-auth-service/models"
	JwtToken "github.com/sourav014/go-auth-service/token"
	"github.com/sourav014/go-auth-service/util"
	"gorm.io/gorm"
)

func GetJSONString(obj interface{}, ignoreFields ...string) (string, error) {
	toJson, err := json.Marshal(obj)
	if err != nil {
		return "", err
	}
	if len(ignoreFields) == 0 {
		return string(toJson), nil
	}
	toMap := map[string]interface{}{}
	json.Unmarshal([]byte(string(toJson)), &toMap)

	for _, field := range ignoreFields {
		delete(toMap, field)
	}

	toJson, err = json.Marshal(toMap)
	if err != nil {
		return "", err
	}

	return string(toJson), nil
}

func RegisterUser(ctx *gin.Context) {
	db, ok := ctx.MustGet("database").(gorm.DB)
	if !ok {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "internal server error",
		})
		return
	}
	fmt.Println("db is ", db)
	var registerUserInput RegisterUserInput
	if err := ctx.ShouldBindJSON(&registerUserInput); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	fmt.Println("registerUserInput is ", registerUserInput)
	var existingUser models.User
	result := db.Where("email = ?", registerUserInput.Email).First(&existingUser)

	if result.Error != nil && !errors.Is(result.Error, gorm.ErrRecordNotFound) {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "internal server error",
		})
		return
	}

	if existingUser.ID != 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "email already exists",
		})
		return
	}

	hashedPassword, err := util.GenerateHashString(registerUserInput.Password)
	if err != nil {
		log.Printf("Error while generating hash string: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "internal server error",
		})
		return
	}

	newUser := models.User{
		Name:     registerUserInput.Name,
		Email:    registerUserInput.Email,
		Password: hashedPassword,
		IsAdmin:  registerUserInput.IsAdmin,
	}

	if err := db.Create(&newUser).Error; err != nil {
		log.Printf("Error creating user: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "internal server error",
		})
		return
	}

	userResp := RegisterUserResp{
		Name:    newUser.Name,
		Email:   newUser.Email,
		IsAdmin: newUser.IsAdmin,
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"user": userResp,
	})
}

func LoginUser(ctx *gin.Context) {
	db, ok := ctx.MustGet("database").(gorm.DB)
	if !ok {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "internal server error",
		})
		return
	}

	var loginUserInput LoginUserInput
	if err := ctx.ShouldBindJSON(&loginUserInput); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	var user models.User
	result := db.Where("email = ?", loginUserInput.Email).First(&user)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{
				"error": "email not found",
			})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "internal server error",
		})
		return
	}

	err := util.CompareHashString(user.Password, loginUserInput.Password)
	if err != nil {
		ctx.JSON(http.StatusForbidden, gin.H{
			"error": "invalid credentials",
		})
		return
	}

	jwtMaker := JwtToken.NewJWTMaker(os.Getenv("SECRET_KEY"))

	accessToken, _, err := jwtMaker.CreateToken(user.ID, user.IsAdmin, time.Minute*15)
	if err != nil {
		log.Printf("Error while creating access token: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to generate token",
		})
		return
	}
	refreshToken, refreshClaims, err := jwtMaker.CreateToken(user.ID, user.IsAdmin, time.Hour*24)
	if err != nil {
		log.Printf("Error while creating refresh token: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to generate token",
		})
		return
	}
	session := models.Session{
		ID:           refreshClaims.RegisteredClaims.ID,
		UserEmail:    user.Email,
		RefreshToken: refreshToken,
		IsRevoked:    false,
		ExpiresAt:    refreshClaims.RegisteredClaims.ExpiresAt.Time,
	}

	if err := db.Create(&session).Error; err != nil {
		log.Printf("Error creating session: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "internal server error",
		})
		return
	}

	loginResp := LoginUserResp{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		SessionId:    session.ID,
		User: RegisterUserResp{
			Name:    user.Name,
			Email:   user.Email,
			IsAdmin: user.IsAdmin,
		},
	}

	ctx.JSON(http.StatusOK, loginResp)
}

func RenewAccessToken(ctx *gin.Context) {
	db, ok := ctx.MustGet("database").(gorm.DB)
	if !ok {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "internal server error",
		})
		return
	}

	var renewAccessTokenInput RenewAccessTokenInput
	if err := ctx.ShouldBindJSON(&renewAccessTokenInput); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	jwtMaker := JwtToken.NewJWTMaker(os.Getenv("SECRET_KEY"))

	refreshClaims, err := jwtMaker.VerifyToken(renewAccessTokenInput.RefreshToken)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	var session models.Session
	result := db.Where("id = ?", refreshClaims.RegisteredClaims.ID).First(&session)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{
				"error": "id not found",
			})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "internal server error",
		})
		return
	}

	if session.IsRevoked {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error": "session has revoked",
		})
		return
	}

	accessToken, accessClaims, err := jwtMaker.CreateToken(refreshClaims.ID, refreshClaims.IsAdmin, time.Minute*15)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "internal server error",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"access_token":            accessToken,
		"access_toker_expires_at": accessClaims.RegisteredClaims.ExpiresAt.Time,
	})
}

func RevokeToken(ctx *gin.Context) {
	db, ok := ctx.MustGet("database").(gorm.DB)
	if !ok {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "internal server error",
		})
		return
	}
	sessionId := ctx.Param("id")
	result := db.Model(&models.Session{}).Where("id = ?", sessionId).Update("is_revoked", true)

	if result.Error != nil {
		log.Printf("Error updating session: %v", result.Error)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "internal server error",
		})
		return
	}

	if result.RowsAffected == 0 {
		log.Printf("No session found with ID: %v", sessionId)
		ctx.JSON(http.StatusNotFound, gin.H{
			"error": "session not found",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "token revoked successfully",
	})
}

func GetUserProfile(ctx *gin.Context) {
	currentUser, ok := ctx.MustGet("currentUser").(models.User)
	if !ok {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "internal server error",
		})
		return
	}

	userProfile := RegisterUserResp{
		Name:    currentUser.Name,
		Email:   currentUser.Email,
		IsAdmin: currentUser.IsAdmin,
	}

	ctx.JSON(http.StatusOK, gin.H{
		"user": userProfile,
	})
}
