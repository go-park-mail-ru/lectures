package repository

import (
	"crudapp/internal/pkg/models"
)

type UserRepo struct {
	data map[string]*models.User
}

func NewUserRepo() *UserRepo {
	return &UserRepo{
		data: map[string]*models.User{
			"rvasily": {
				ID:       1,
				Login:    "rvasily",
				Password: "love",
			},
		},
	}
}

func (repo *UserRepo) Authorize(login, pass string) (*models.User, error) {
	u, ok := repo.data[login]
	if !ok {
		return nil, models.ErrNoUser
	}

	// dont do this un production :)
	if u.Password != pass {
		return nil, models.ErrBadPass
	}

	return u, nil
}
