package main

import (
	"log"
	"os"

	"be-go-test-thai-bev-auth/config"
	"be-go-test-thai-bev-auth/internal/handler"
	"be-go-test-thai-bev-auth/internal/repository"
	"be-go-test-thai-bev-auth/internal/router"
	"be-go-test-thai-bev-auth/internal/service"

	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()

	db, err := config.ConnectDB()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	userRepo := repository.NewUserRepository(db)
	authSvc := service.NewAuthService(userRepo)
	authHandler := handler.NewAuthHandler(authSvc)

	r := router.Setup(authHandler)

	port := os.Getenv("SERVER_PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server running on port %s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
