package main

import (
	"fmt"
	"log"
	"os"

	"github.com/go-playground/validator/v10"
	"github.com/sourav014/go-auth-service/controller"
	"github.com/sourav014/go-auth-service/db"
	"github.com/sourav014/go-auth-service/middleware"
	"github.com/sourav014/go-auth-service/models"
	"github.com/sourav014/go-auth-service/repository"
	"github.com/sourav014/go-auth-service/router"
	"github.com/sourav014/go-auth-service/service"
	JwtToken "github.com/sourav014/go-auth-service/token"
)

func main() {
	log.Println("Started Server!")
	database, err := db.NewDatabase()
	if err != nil {
		fmt.Println("error while creating database instance", err)
	}

	if err := database.GetDB().AutoMigrate(&models.User{}); err != nil {
		fmt.Println("Error during user migration:", err)
		return
	}

	if err := database.GetDB().AutoMigrate(&models.Session{}); err != nil {
		fmt.Println("Error during session migration:", err)
		return
	}

	usersRepository := repository.NewUsersRepositoryImpl(database.GetDB())
	sessionsRepository := repository.NewSessionsRepositoryImpl(database.GetDB())

	jwtMaker := JwtToken.NewJWTMaker(os.Getenv("SECRET_KEY"))
	validate := validator.New()

	authService := service.NewAuthServiceImpl(sessionsRepository, usersRepository, jwtMaker, validate)
	authController := controller.NewAuthController(authService)

	authMiddleware := middleware.NewAuthMiddlewareImpl(sessionsRepository, usersRepository, jwtMaker, validate)

	router := router.NewRouter(authController, authMiddleware)

	router.Run()
}
