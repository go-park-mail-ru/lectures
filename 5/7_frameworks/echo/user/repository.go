package user

import (
	"echo_example/model"
	"errors"
)

var ErrUserExists = errors.New("user exists")

type Repository interface {
	GetUser(username string) (model.User, error)
	GetAllUsers() ([]model.User, error)

	InsertUser(user model.User) error
}
