package repository

import (
	"errors"

	"github.com/sourav014/go-auth-service/helper"
	"github.com/sourav014/go-auth-service/models"
	"gorm.io/gorm"
)

type SessionsRepositoryImpl struct {
	Db *gorm.DB
}

func NewSessionsRepositoryImpl(Db *gorm.DB) SessionsRepository {
	return &SessionsRepositoryImpl{Db: Db}
}

func (s *SessionsRepositoryImpl) Create(session models.Session) {
	result := s.Db.Create(&session)
	helper.ErrorPanic(result.Error)
}

func (s *SessionsRepositoryImpl) Update(id string) {
	result := s.Db.Model(&models.Session{}).Where("id = ?", id).Update("is_revoked", true)
	helper.ErrorPanic(result.Error)
}

func (s *SessionsRepositoryImpl) FindById(id string) (models.Session, error) {
	var session models.Session
	result := s.Db.Where("id = ?", id).First(&session)
	if result != nil {
		return session, nil
	} else {
		return session, errors.New("session not found")
	}
}
