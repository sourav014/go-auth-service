package models

import (
	"time"
)

type User struct {
	ID        uint   `gorm:"primaryKey;autoIncrement;not null"`
	Name      string `gorm:"not null"`
	Email     string `gorm:"unique;not null"`
	Password  string `gorm:"not null"`
	IsAdmin   bool
	CreatedAt time.Time
	UpdatedAt time.Time
}

type Session struct {
	ID           string `gorm:"primaryKey;not null"`
	UserEmail    string `gorm:"not null"`
	RefreshToken string `gorm:"not null"`
	IsRevoked    bool
	CreatedAt    time.Time
	ExpiresAt    time.Time
}
