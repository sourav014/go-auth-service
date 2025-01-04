package middleware

import "github.com/gin-gonic/gin"

type AuthMiddleware interface {
	CheckUserAuthentication(ctx *gin.Context)
}
