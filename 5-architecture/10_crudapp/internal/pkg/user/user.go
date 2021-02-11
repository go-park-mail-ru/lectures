package user

import "crudapp/internal/pkg/models"

//go:generate mockgen -destination=./mock/mock_repo.go -package=mock crudapp/internal/pkg/user Repository

type Repository interface {
	GetByLogin(login string) (*models.User, error)
	Authorize(login, pass string) (*models.User, error)
}
