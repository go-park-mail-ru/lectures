package items

import "crudapp/internal/pkg/models"

type Repository interface {
	GetAll() ([]*models.Item, error)
	GetByID(id uint32) (*models.Item, error)
	Add(item *models.Item) (uint32, error)
	Update(newItem *models.Item) (bool, error)
	Delete(id uint32) (bool, error)
}
