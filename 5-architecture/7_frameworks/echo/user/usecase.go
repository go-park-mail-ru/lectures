package user

import "github.com/go-park-mail-ru/lectures/5-architecture/7_frameworks/echo/model"

type Usecase interface {
	GetUser(username string) (model.User, error)
	GetAllUsers() ([]model.User, error)

	CreateUser(user model.User) error
}
