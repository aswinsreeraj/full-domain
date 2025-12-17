package domain

import (
	"context"
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

type UserRepository interface {
	Create(user *User) error
	FindByEmail(email string) (*User, error)
	FindByID(id uint) (*User, error)
	Update(user *User) error
	Delete(id uint) error
	Search(query string) ([]User, error)
	FindAll() ([]User, error)
}

type UserService interface {
	CreateUser(name, email, password string) error
	Authenticate(ctx context.Context, email, password string) (*User, error)
	Update(user *User) error
	UpdatePassword(email, oldPassword, newPassword string) error
	SearchUsers(query string) ([]User, error)
	FindByIDString(id string) (*User, error)
	UpdateUser(ctx context.Context, id, name, email, role, password string) error
	DeleteUser(id string) error
	FindByEmail(email string) (*User, error)
}
