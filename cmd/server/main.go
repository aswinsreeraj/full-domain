package main

import (
	"log"
	"os"

	httpDelivery "full-domain/internal/delivery/http"
	"full-domain/internal/repository/postgres"
	"full-domain/internal/usecase"
	"full-domain/pkg/woodpecker"

	"github.com/joho/godotenv"
)

func main() {

	logFile, err := os.OpenFile("app.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic("failed to create log file: " + err.Error())
	}
	defer logFile.Close()

	woodpecker.Init(logFile)
	woodpecker.Logger.Info("starting application")

	woodpecker.Logger.Info("Loading env file")
	if err := godotenv.Load(); err != nil {
		woodpecker.Logger.Error("failed to load env", "error", err)
	}

	if err := postgres.Connect(); err != nil {
		woodpecker.Logger.Error("database connection failed", "error", err)
	}

	userRepo := postgres.NewUserRepository(postgres.DB)
	userService := usecase.NewUserService(userRepo)

	r := httpDelivery.NewRouter(userService)

	log.Fatal(r.Run(":8080"))
}
