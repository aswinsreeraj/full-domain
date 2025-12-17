package postgres

import (
	"full-domain/pkg/woodpecker"
	"full-domain/internal/domain"
	"log/slog"
	"os"

	"gorm.io/gorm"
)



type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) domain.UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(user *domain.User) error {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	logger.Info("Creating user", "email", user.Email)
	err := r.db.Create(user).Error
	if err != nil {
		woodpecker.Logger.Error("failed to create user", "email", user.Email, "error", err)
	}
	return err
	// Create is better than Save for new records
	// as it only inserts and does not check for existing records
}

func (r *userRepository) FindByEmail(email string) (*domain.User, error) {
	var user domain.User
	err := r.db.Where("email = ?", email).First(&user).Error
	// If no record found, GORM returns gorm.ErrRecordNotFound
	return &user, err
}

func (r *userRepository) FindByID(id uint) (*domain.User, error) {
	var user domain.User
	err := r.db.First(&user, id).Error
	// If id == 0 or no record found, GORM returns gorm.ErrRecordNotFound
	return &user, err
}

func (r *userRepository) Update(user *domain.User) error {
	return r.db.Model(&domain.User{}).Where("id = ?", user.ID).Updates(user).Error
	// Updates ignores zero values, unlike Save
	// r.db.Save(user).Error - earlier approach
	// If zero values need to be updated, use Select("*").Updates(user) or use map[string]interface{}
	// Save - PUT, whereas Updates - PATCH
}

func (r *userRepository) Delete(id uint) error {
	return r.db.Delete(&domain.User{}, id).Error
	// Soft delete
	// r.db.Unscoped().Delete(&models.User{}, id).Error - hard delete
	// Hard delete should be reserved for admin operations only
}

func (r *userRepository) Search(query string) ([]domain.User, error) {
	var users []domain.User
	err := r.db.Where("LOWER(name) LIKE LOWER(?) OR LOWER(email) LIKE LOWER(?)", "%"+query+"%", "%"+query+"%").Find(&users).Error
	// LOWER for case-insensitive search
	return users, err
}

func (r *userRepository) FindAll() ([]domain.User, error) {
	var users []domain.User
	err := r.db.Find(&users).Error
	return users, err
}

func (r *userRepository) GetDB() *gorm.DB {
	return r.db
}
