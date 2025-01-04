package controller

import (
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sourav014/go-auth-service/initializers"
	"github.com/sourav014/go-auth-service/models"
	JwtToken "github.com/sourav014/go-auth-service/token"
)

func CheckUserAuthentication(ctx *gin.Context) {
	authHeader := ctx.GetHeader("Authorization")

	if authHeader == "" {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is missing"})
		return
	}

	authToken := strings.Split(authHeader, " ")
	if len(authToken) != 2 || authToken[0] != "Bearer" {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token format"})
		return
	}

	tokenString := authToken[1]
	jwtMaker := JwtToken.NewJWTMaker(os.Getenv("SECRET_KEY"))

	userClaims, err := jwtMaker.VerifyToken(tokenString)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
		return
	}

	if time.Now().Unix() > userClaims.ExpiresAt.Unix() {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Token expired"})
		return
	}

	var user models.User
	initializers.DB.Where("ID=?", userClaims.ID).Find(&user)

	if user.ID == 0 {
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	ctx.Set("currentUser", user)

	ctx.Next()
}
