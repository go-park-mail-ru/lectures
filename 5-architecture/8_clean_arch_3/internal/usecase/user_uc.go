// Package usecase реализует сценарии приложения и задаёт контракт хранилища.
package usecase

import (
	"context"
	"errors"
	"fmt"

	"lectures/cleanarch3/internal/domain"
)

var (
	ErrNotFound     = errors.New("user not found")
	ErrEmailTaken   = errors.New("email already taken")
	ErrInvalidInput = errors.New("invalid input")
)

// UserRepository объявлен у потребителя — сценариев, а не у реализации БД.
type UserRepository interface {
	Create(context.Context, domain.User) (domain.User, error)
	GetByID(context.Context, int64) (domain.User, error)
	List(context.Context, int64, int) ([]domain.User, error)
	Update(context.Context, domain.User) error
	Delete(context.Context, int64) error
}

// Модели сценариев не содержат JSON-тегов или HTTP-типов.
type CreateUserInput struct {
	Name, Email string
}

type UpdateUserInput struct {
	ID          int64
	Name, Email string
}

type UserPage struct {
	Items  []domain.User
	NextID int64
}

type UserUseCase struct {
	users UserRepository
}

func NewUserUseCase(users UserRepository) *UserUseCase {
	return &UserUseCase{users: users}
}

func (uc *UserUseCase) Register(ctx context.Context, in CreateUserInput) (domain.User, error) {
	u, err := domain.NewUser(in.Name, in.Email)
	if err != nil {
		return domain.User{}, err
	}
	// Поиск email перед Create не защитил бы от двух одновременных запросов.
	// Уникальность обеспечивает репозиторий при записи.
	u, err = uc.users.Create(ctx, u)
	if err != nil {
		return domain.User{}, fmt.Errorf("register user: %w", err)
	}
	return u, nil
}

func (uc *UserUseCase) GetByID(ctx context.Context, id int64) (domain.User, error) {
	if id <= 0 {
		return domain.User{}, ErrInvalidInput
	}
	return uc.users.GetByID(ctx, id)
}

func (uc *UserUseCase) List(ctx context.Context, afterID int64, limit int) (UserPage, error) {
	if afterID < 0 || limit < 0 || limit > 100 {
		return UserPage{}, ErrInvalidInput
	}
	if limit == 0 {
		limit = 10
	}
	items, err := uc.users.List(ctx, afterID, limit+1)
	if err != nil {
		return UserPage{}, err
	}
	page := UserPage{Items: items}
	if len(items) > limit {
		page.Items = items[:limit]
		page.NextID = page.Items[limit-1].ID()
	}
	return page, nil
}

func (uc *UserUseCase) Update(ctx context.Context, in UpdateUserInput) (domain.User, error) {
	if in.ID <= 0 {
		return domain.User{}, ErrInvalidInput
	}
	u, err := uc.users.GetByID(ctx, in.ID)
	if err != nil {
		return domain.User{}, err
	}
	u, err = u.ChangeProfile(in.Name, in.Email)
	if err != nil {
		return domain.User{}, err
	}
	if err := uc.users.Update(ctx, u); err != nil {
		return domain.User{}, fmt.Errorf("update user: %w", err)
	}
	return u, nil
}

func (uc *UserUseCase) Delete(ctx context.Context, id int64) error {
	if id <= 0 {
		return ErrInvalidInput
	}
	return uc.users.Delete(ctx, id)
}
