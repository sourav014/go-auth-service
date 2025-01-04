package repository

import "github.com/sourav014/go-auth-service/models"

type UsersRepository interface {
	Create(user models.User)
	FindByEmail(email string) (models.User, error)
	FindById(id uint) (models.User, error)
}
