package user

import "crudapp/internal/pkg/models"

type Repository interface {
	Authorize(login, pass string) (*models.User, error)
}
