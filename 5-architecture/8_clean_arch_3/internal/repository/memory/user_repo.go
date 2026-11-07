// Package memory — адаптер хранилища для запуска без PostgreSQL.
package memory

import (
	"context"
	"sort"
	"sync"

	"lectures/cleanarch3/internal/domain"
	"lectures/cleanarch3/internal/usecase"
)

type UserRepository struct {
	mu      sync.RWMutex
	lastID  int64
	users   map[int64]domain.User
	byEmail map[string]int64
}

func NewUserRepository() *UserRepository {
	return &UserRepository{
		users:   make(map[int64]domain.User),
		byEmail: make(map[string]int64),
	}
}

var _ usecase.UserRepository = (*UserRepository)(nil)

func (r *UserRepository) Create(ctx context.Context, u domain.User) (domain.User, error) {
	if err := ctx.Err(); err != nil {
		return domain.User{}, err
	}
	if u.ID() != 0 {
		return domain.User{}, usecase.ErrInvalidInput
	}
	// Проверяем и нулевое значение User{}, которое можно создать в Go.
	if _, err := domain.NewUser(u.Name(), u.Email()); err != nil {
		return domain.User{}, err
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, exists := r.byEmail[u.Email()]; exists {
		return domain.User{}, usecase.ErrEmailTaken
	}
	created, err := domain.RestoreUser(r.lastID+1, u.Name(), u.Email())
	if err != nil {
		return domain.User{}, err
	}
	r.lastID = created.ID()
	r.users[created.ID()] = created
	r.byEmail[created.Email()] = created.ID()
	return created, nil
}

func (r *UserRepository) GetByID(ctx context.Context, id int64) (domain.User, error) {
	if err := ctx.Err(); err != nil {
		return domain.User{}, err
	}
	r.mu.RLock()
	defer r.mu.RUnlock()
	u, exists := r.users[id]
	if !exists {
		return domain.User{}, usecase.ErrNotFound
	}
	return u, nil
}

func (r *UserRepository) List(ctx context.Context, afterID int64, limit int) ([]domain.User, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	if afterID < 0 || limit <= 0 {
		return nil, usecase.ErrInvalidInput
	}
	r.mu.RLock()
	defer r.mu.RUnlock()
	items := make([]domain.User, 0)
	for id, u := range r.users {
		if id > afterID {
			items = append(items, u)
		}
	}
	sort.Slice(items, func(i, j int) bool { return items[i].ID() < items[j].ID() })
	if len(items) > limit {
		items = items[:limit]
	}
	return items, nil
}

func (r *UserRepository) Update(ctx context.Context, u domain.User) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	if _, err := domain.RestoreUser(u.ID(), u.Name(), u.Email()); err != nil {
		return err
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	previous, exists := r.users[u.ID()]
	if !exists {
		return usecase.ErrNotFound
	}
	if id, exists := r.byEmail[u.Email()]; exists && id != u.ID() {
		return usecase.ErrEmailTaken
	}
	delete(r.byEmail, previous.Email())
	r.users[u.ID()] = u
	r.byEmail[u.Email()] = u.ID()
	return nil
}

func (r *UserRepository) Delete(ctx context.Context, id int64) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	u, exists := r.users[id]
	if !exists {
		return usecase.ErrNotFound
	}
	delete(r.byEmail, u.Email())
	delete(r.users, id)
	return nil
}
