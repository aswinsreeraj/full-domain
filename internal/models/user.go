package models

import "time"

const (
	RoleUser  = "user"
	RoleAdmin = "admin"
)

type User struct {
	ID             uint   `gorm:"primaryKey"`
	Name           string `gorm:"type:varchar(100);not null"`
	Email          string `gorm:"type:varchar(255);unique;not null"`
	HashedPassword string `gorm:"type:varchar(60);not null"`
	Role           string `gorm:"type:user_role;not null"`
	LastLogin      *time.Time
	CreatedAt      time.Time `gorm:"not null;default;current_timestamp"`
	UpdatedAt      time.Time `gorm:"not null;default;current_timestamp"`
}
