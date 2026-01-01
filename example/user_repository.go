package example

import (
	"github.com/mhcvintar/repo-gen/repository"
)

type UserRepository interface {
	repository.Repository[repository.UserModel, int64]

	FindByEmail(email string) (*repository.UserModel, error)
	FindByFirstNameAndLastName(firstName, lastName string) ([]*repository.UserModel, error)
}
