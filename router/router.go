package router

import (
	"github.com/gin-gonic/gin"
	"github.com/sourav014/go-auth-service/controller"
	"github.com/sourav014/go-auth-service/middleware"
)

func healthCheck(ctx *gin.Context) {
	ctx.JSON(200, gin.H{
		"message": "server is up..",
	})
}

func NewRouter(authController *controller.AuthController, authMiddleware middleware.AuthMiddleware) *gin.Engine {
	router := gin.Default()
	api := router.Group("/api")
	{
		v1 := api.Group("v1")
		{
			auth := v1.Group("/auth")
			{
				auth.POST("/register", authController.RegisterUser)
				auth.POST("/login", authController.LoginUser)
				auth.POST("/renew", authController.RenewAccessToken)
				auth.POST("/revoke/:id", authController.RevokeToken)
			}
			user := v1.Group("/user")
			{
				user.GET("/profile", authMiddleware.CheckUserAuthentication, authController.GetUserProfile)
			}
			health := v1.Group("/health")
			{
				health.GET("/check", healthCheck)
			}
		}

	}

	return router
}
