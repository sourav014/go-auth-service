package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/sourav014/go-auth-service/controller"
	"github.com/sourav014/go-auth-service/db"
	"github.com/sourav014/go-auth-service/handler"
	"github.com/sourav014/go-auth-service/initializers"
	"gorm.io/gorm"
)

func healthCheck(ctx *gin.Context) {
	ctx.JSON(200, gin.H{
		"message": "server is up..",
	})
}

func ApiMiddleware(db gorm.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Set("database", db)
		ctx.Next()
	}
}

func initRouter(db *gorm.DB) *gin.Engine {
	router := gin.Default()
	router.Use(ApiMiddleware(*db))

	api := router.Group("/api")
	{
		v1 := api.Group("v1")
		{
			auth := v1.Group("/auth")
			{
				auth.POST("/register", controller.RegisterUser)
				auth.POST("/login", controller.LoginUser)
				auth.POST("/renew", controller.RenewAccessToken)
				auth.POST("/revoke/:id", controller.RevokeToken)
			}
			user := v1.Group("/user")
			{
				user.GET("/profile", controller.CheckUserAuthentication, controller.GetUserProfile)
			}
			health := v1.Group("/health")
			{
				health.GET("/check", healthCheck)
			}
		}

	}

	return router
}

func init() {
	initializers.LoadEnvs()
	initializers.ConnectDB()
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	db, err := db.NewDatabase()
	if err != nil {
		fmt.Println("error while creating database instance", err)
	}
	handler.NewHandler(db.GetDB(), os.Getenv("SECRET_KEY"))
	router := initRouter(initializers.DB)
	router.Run()
}
