package repository

import (
	"errors"

	"github.com/sourav014/go-auth-service/helper"
	"github.com/sourav014/go-auth-service/models"
	"gorm.io/gorm"
)

type UsersRepositoryImpl struct {
	Db *gorm.DB
}

func NewUsersRepositoryImpl(Db *gorm.DB) UsersRepository {
	return &UsersRepositoryImpl{Db: Db}
}

func (u *UsersRepositoryImpl) Create(user models.User) {
	result := u.Db.Create(&user)
	helper.ErrorPanic(result.Error)
}

func (u *UsersRepositoryImpl) FindByEmail(email string) (models.User, error) {
	var user models.User
	result := u.Db.Where("email = ?", email).First(&user)
	if result != nil {
		return user, nil
	} else {
		return user, errors.New("user not found")
	}
}

func (u *UsersRepositoryImpl) FindById(id uint) (models.User, error) {
	var user models.User
	result := u.Db.Where("id = ?", id).First(&user)
	if result != nil {
		return user, nil
	} else {
		return user, errors.New("user not found")
	}
}
