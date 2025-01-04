package repository

import "github.com/sourav014/go-auth-service/models"

type SessionsRepository interface {
	Create(session models.Session)
	Update(id string)
	FindById(id string) (models.Session, error)
}
