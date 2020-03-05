package items

// WARNING! completly unsafe in multi-goroutine use, need add mutexes

//go:generate mockgen -destination=./mock_repo.go -package=items crudapp/pkg/items ItemsRepo

type ItemsRepo interface {
	GetAll() ([]*Item, error)
	GetByID(id uint32) (*Item, error)
	Add(item *Item) (uint32, error)
	Update(newItem *Item) (bool, error)
	Delete(id uint32) (bool, error)
}

type itemsRepo struct {
	lastID uint32
	data   []*Item
}

func NewRepo() *itemsRepo {
	return &itemsRepo{
		data: make([]*Item, 0, 10),
	}
}

func (repo *itemsRepo) GetAll() ([]*Item, error) {
	return repo.data, nil
}

func (repo *itemsRepo) GetByID(id uint32) (*Item, error) {
	for _, item := range repo.data {
		if item.ID == id {
			return item, nil
		}
	}
	return nil, nil
}

func (repo *itemsRepo) Add(item *Item) (uint32, error) {
	repo.lastID++
	item.ID = repo.lastID
	repo.data = append(repo.data, item)
	return repo.lastID, nil
}

func (repo *itemsRepo) Update(newItem *Item) (bool, error) {
	for _, item := range repo.data {
		if item.ID != newItem.ID {
			continue
		}
		item.Title = newItem.Title
		item.Description = newItem.Description
		return true, nil
	}
	return false, nil
}

func (repo *itemsRepo) Delete(id uint32) (bool, error) {
	i := -1
	for idx, item := range repo.data {
		if item.ID != id {
			continue
		}
		i = idx
	}
	if i < 0 {
		return false, nil
	}

	if i < len(repo.data)-1 {
		copy(repo.data[i:], repo.data[i+1:])
	}
	repo.data[len(repo.data)-1] = nil // or the zero value of T
	repo.data = repo.data[:len(repo.data)-1]

	return true, nil
}
