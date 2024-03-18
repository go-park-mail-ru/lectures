package items

import (
	"log"

	"gorm.io/gorm"
)

type RepoGorm struct {
	DB *gorm.DB
}

func NewGormRepository(db *gorm.DB) *RepoGorm {
	return &RepoGorm{DB: db}
}

func (repo *RepoGorm) GetAll() ([]*Item, error) {
	// https://gorm.io/docs/query.html
	items := make([]*Item, 0, 10)
	err := repo.DB.Find(&items).Error
	if err != nil {
		return nil, err
	}
	return items, nil
}

func (repo *RepoGorm) GetByID(id int64) (*Item, error) {
	post := &Item{}
	err := repo.DB.First(&post, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return post, nil
}

func (repo *RepoGorm) Add(elem *Item) (int64, error) {
	err := repo.DB.Create(elem).Error
	log.Println("created elem id:", elem.ID)
	if err != nil {
		return 0, err
	}
	return int64(elem.ID), nil
}

func (repo *RepoGorm) Update(elem *Item) (int64, error) {
	res := repo.DB.Model(&elem).Updates(map[string]interface{}{
		"title":       elem.Title,
		"description": elem.Description,
		"updated":     "rvasily",
	})
	if res.Error != nil {
		return 0, res.Error
	}
	return res.RowsAffected, nil
}

func (repo *RepoGorm) Delete(id int64) (int64, error) {
	res := repo.DB.Delete(&Item{}, id)
	if res.Error != nil {
		return 0, res.Error
	}
	return res.RowsAffected, nil
}
