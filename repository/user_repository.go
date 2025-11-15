package repository

import (
	"github.com/mhcvintar/repo-gen/models"
)

//go:generate go run ../cmd/generator.go -model=User -input=user_repository.go

type UserRepository interface {
	FindByName(name string) (*models.User, error)
	FindByEmail(email string) (*models.User, error)
}