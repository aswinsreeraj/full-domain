package services

import (
	"full-domain/internal/database"
	"full-domain/internal/lumberjack"
	"full-domain/internal/models"
	"strconv"

	"golang.org/x/crypto/bcrypt"
)

type UserService interface {
	CreateUser(name, email, password string) error
	Authenticate(email, password string) (*models.User, error)
	Update(user *models.User) error
	UpdatePassword(email, oldPassword, newPassword string) error
	SearchUsers(query string) ([]models.User, error)
	FindByIDString(id string) (*models.User, error)
	UpdateUser(id, name, email, role, password string) error
	DeleteUser(id string) error
	FindByEmail(email string) (*models.User, error)
}

type userService struct {
	repo database.UserRepository
}

func NewUserService(repo database.UserRepository) UserService {
	return &userService{repo: repo}
}

func (s *userService) CreateUser(name, email, password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user := &models.User{
		Name:           name,
		Email:          email,
		HashedPassword: string(hashedPassword),
		Role:           models.RoleUser,
	}

	if err := s.repo.Create(user); err != nil {
		return err
	}
	return nil
}

func (s *userService) Authenticate(email, password string) (*models.User, error) {
	lumberjack.Logger.Info("authenticating user", "email", email)
	user, err := s.repo.FindByEmail(email)
	if err != nil {
		lumberjack.Logger.Warn("authentication failed: user not found", "email", email)
		return nil, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.HashedPassword), []byte(password)); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *userService) Update(user *models.User) error {
	return s.repo.Update(user)
}

func (s *userService) UpdatePassword(email, oldPassword, newPassword string) error {
	user, err := s.repo.FindByEmail(email)
	if err != nil {
		return err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.HashedPassword), []byte(oldPassword)); err != nil {
		return err
	}

	hashedNew, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user.HashedPassword = string(hashedNew)
	return s.repo.Update(user)
}

func (s *userService) SearchUsers(query string) ([]models.User, error) {
	if query == "" {
		return s.repo.FindAll()
	}
	return s.repo.Search(query)
}

func (s *userService) FindByIDString(id string) (*models.User, error) {
	uid, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		return nil, err
	}
	return s.repo.FindByID(uint(uid))
}

func (s *userService) UpdateUser(id, name, email, role, password string) error {
	lumberjack.Logger.Info("admin updating user", "id", id, "role", role)
	uid, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		return err
	}

	user, err := s.repo.FindByID(uint(uid))
	if err != nil {
		return err
	}

	user.Name = name
	user.Email = email
	user.Role = role
	if password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			return err
		}
		user.HashedPassword = string(hashedPassword)
	}
	return s.repo.Update(user)
}

func (s *userService) DeleteUser(id string) error {
	uid, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		return err
	}
	return s.repo.Delete(uint(uid))
}

func (s *userService) FindByEmail(email string) (*models.User, error) {
	return s.repo.FindByEmail(email)
}
