package handler

import (
	"github.com/sourav014/go-auth-service/token"
	"gorm.io/gorm"
)

type handler struct {
	db         *gorm.DB
	TokenMaker *token.JWTMaker
}

func NewHandler(db *gorm.DB, secretKey string) *handler {
	return &handler{
		db:         db,
		TokenMaker: token.NewJWTMaker(secretKey),
	}
}
