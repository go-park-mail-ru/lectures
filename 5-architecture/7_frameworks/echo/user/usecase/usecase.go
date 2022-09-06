package usecase

import (
	"github.com/go-park-mail-ru/lectures/5-architecture/7_frameworks/echo/model"
	"github.com/go-park-mail-ru/lectures/5-architecture/7_frameworks/echo/user"
)

func NewUserUsecase(userRepo user.Repository) user.Usecase {
	return userUsecase{repo: userRepo}
}

type userUsecase struct {
	repo user.Repository
}

func (u userUsecase) GetUser(username string) (model.User, error) {
	return u.repo.GetUser(username)
}

func (u userUsecase) GetAllUsers() ([]model.User, error) {
	return u.repo.GetAllUsers()
}

func (u userUsecase) CreateUser(user model.User) error {
	return u.repo.InsertUser(user)
}
