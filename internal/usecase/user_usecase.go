package usecase

import (
	"context"
	"full-domain/internal/domain"
	"full-domain/pkg/woodpecker"
	"strconv"

	"golang.org/x/crypto/bcrypt"
)

type userService struct {
	repo domain.UserRepository
}

func NewUserService(repo domain.UserRepository) domain.UserService {
	return &userService{repo: repo}
}

func (s *userService) CreateUser(name, email, password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user := &domain.User{
		Name:           name,
		Email:          email,
		HashedPassword: string(hashedPassword),
		Role:           domain.RoleUser,
	}

	if err := s.repo.Create(user); err != nil {
		return err
	}
	return nil
}

func (s *userService) Authenticate(ctx context.Context, email, password string) (*domain.User, error) {
	logger := woodpecker.FromContext(ctx)
	logger.Info("authenticating user", "email", email)

	user, err := s.repo.FindByEmail(email)
	if err != nil {
		logger.Warn("authentication failed: user not found", "email", email)
		return nil, err
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.HashedPassword), []byte(password)); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *userService) Update(user *domain.User) error {
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

func (s *userService) SearchUsers(query string) ([]domain.User, error) {
	if query == "" {
		return s.repo.FindAll()
	}
	return s.repo.Search(query)
}

func (s *userService) FindByIDString(id string) (*domain.User, error) {
	uid, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		return nil, err
	}
	return s.repo.FindByID(uint(uid))
}

func (s *userService) UpdateUser(ctx context.Context, id, name, email, role, password string) error {
	logger := woodpecker.FromContext(ctx)
	logger.Info("admin updating user", "id", id, "role", role)
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

func (s *userService) FindByEmail(email string) (*domain.User, error) {
	return s.repo.FindByEmail(email)
}
