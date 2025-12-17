package postgres

import (
	"full-domain/pkg/woodpecker"
	"full-domain/internal/domain"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Connect() error {
	dsn := os.Getenv("DATABASE_DSN")
	woodpecker.Logger.Info("connecting to database")
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		woodpecker.Logger.Error("database connection error", "error", err)
		return err
	}
	DB = db

	if err := db.AutoMigrate(&domain.User{}); err != nil {
		return err
	}
	// Added auto migration for User model
	// Initially added for DeletedAt field
	// Further addition of fields will be handled automatically by GORM

	return nil
}
