package user

import "echo_example/model"

type Usecase interface {
	GetUser(username string) (model.User, error)
	GetAllUsers() ([]model.User, error)

	CreateUser(user model.User) error
}
