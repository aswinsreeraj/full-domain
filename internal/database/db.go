package database

import (
	"full-domain/internal/lumberjack"
	"full-domain/internal/models"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Connect() error {
	dsn := os.Getenv("DATABASE_DSN")
	lumberjack.Logger.Info("connecting to database")
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		lumberjack.Logger.Error("database connection error", "error", err)
		return err
	}
	DB = db

	if err := DB.AutoMigrate(&models.User{}); err != nil {
		return err
	}
	// Added auto migration for User model
	// Initially added for DeletedAt field
	// Further addition of fields will be handled automatically by GORM

	return nil
}
