package main

import (
	"log"

	"github.com/mhcvintar/repo-gen/config"
	"github.com/mhcvintar/repo-gen/models"
	"github.com/mhcvintar/repo-gen/repository"
)

func main() {
	db, err := config.InitDB()
    if err != nil {
        log.Fatalf("Failed to initialize database: %v", err)
    }

	userRepository := repository.NewUserRepository(db)

    user := &models.User{
        Name:  "John Doe",
        Email: "john@example.com",
    }

	result := db.Create(user)
    if result.Error != nil {
        log.Fatalf("Failed to create user: %v", result.Error)
    }
    log.Printf("Created user with ID: %d", user.ID)

	user, err = userRepository.FindByEmail(user.Email)
    if err != nil {
        log.Fatalf("Failed to find user: %v", err)
    }
    log.Printf("Found user with email: %s", user.Email)
}