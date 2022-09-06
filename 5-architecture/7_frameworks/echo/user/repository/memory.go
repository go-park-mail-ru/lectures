package repository

import (
	"sync"

	"github.com/go-park-mail-ru/lectures/5-architecture/7_frameworks/echo/model"
	"github.com/go-park-mail-ru/lectures/5-architecture/7_frameworks/echo/user"
)

func NewUserMemoryRepository() user.Repository {
	return &UserMemoryRepository{
		storage: map[string]model.User{},
	}
}

type UserMemoryRepository struct {
	mu      sync.RWMutex
	storage map[string]model.User
}

func (db *UserMemoryRepository) GetUser(username string) (model.User, error) {
	panic("implement me")
}

func (db *UserMemoryRepository) GetAllUsers() ([]model.User, error) {
	db.mu.RLock()
	defer db.mu.RUnlock()

	ret := make([]model.User, 0, len(db.storage))

	for _, u := range db.storage {
		ret = append(ret, u)
	}

	return ret, nil
}

func (db *UserMemoryRepository) InsertUser(u model.User) error {
	db.mu.Lock()
	defer db.mu.Unlock()

	if _, ok := db.storage[u.Username]; ok {
		return user.ErrUserExists
	}

	db.storage[u.Username] = u

	return nil
}
