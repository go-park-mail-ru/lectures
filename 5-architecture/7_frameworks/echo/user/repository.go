package user

import (
	"errors"

	"github.com/go-park-mail-ru/lectures/5-architecture/7_frameworks/echo/model"
)

var ErrUserExists = errors.New("user exists")

type Repository interface {
	GetUser(username string) (model.User, error)
	GetAllUsers() ([]model.User, error)

	InsertUser(user model.User) error
}
