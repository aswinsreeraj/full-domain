package models

import (
	"time"

	"gorm.io/gorm"
)

const (
	RoleUser  = "user"
	RoleAdmin = "admin"
)

type User struct {
	gorm.Model
	// Includes ID, CreatedAt, UpdatedAt, DeletedAt
	Name           string     `gorm:"type:varchar(100);not null"`
	Email          string     `gorm:"type:varchar(255);unique;not null"`
	HashedPassword string     `gorm:"type:varchar(60);not null"`
	Role           string     `gorm:"type:user_role;not null"`
	LastLogin      *time.Time // Pointer to allow null values
}
