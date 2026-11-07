package usecase_test

import (
	"context"
	"errors"
	"testing"

	"lectures/cleanarch3/internal/domain"
	"lectures/cleanarch3/internal/usecase"
)

// Подменяем только метод, который нужен сценарию регистрации.
// HTTP и реализации хранилищ в этих тестах не участвуют.
type registrationRepository struct {
	usecase.UserRepository
	calls int
	err   error
}

func (r *registrationRepository) Create(_ context.Context, u domain.User) (domain.User, error) {
	r.calls++
	if r.err != nil {
		return domain.User{}, r.err
	}
	return domain.RestoreUser(42, u.Name(), u.Email())
}

func TestRegisterValidatesBeforeStorage(t *testing.T) {
	r := &registrationRepository{}
	uc := usecase.NewUserUseCase(r)
	_, err := uc.Register(context.Background(), usecase.CreateUserInput{Name: " ", Email: "anna@example.com"})
	if !errors.Is(err, domain.ErrInvalidUser) || r.calls != 0 {
		t.Fatalf("invalid user reached storage: calls=%d, err=%v", r.calls, err)
	}
}

func TestRegisterReturnsAssignedIDAndNormalizedProfile(t *testing.T) {
	r := &registrationRepository{}
	uc := usecase.NewUserUseCase(r)
	u, err := uc.Register(context.Background(), usecase.CreateUserInput{Name: " Анна ", Email: "ANNA@example.com"})
	if err != nil || r.calls != 1 || u.ID() != 42 || u.Name() != "Анна" || u.Email() != "anna@example.com" {
		t.Fatalf("registered user: %v, calls=%d, err=%v", u, r.calls, err)
	}
}

func TestRegisterPreservesRepositoryErrors(t *testing.T) {
	for _, original := range []error{usecase.ErrEmailTaken, context.DeadlineExceeded, errors.New("storage unavailable")} {
		r := &registrationRepository{err: original}
		uc := usecase.NewUserUseCase(r)
		_, err := uc.Register(context.Background(), usecase.CreateUserInput{Name: "Анна", Email: "anna@example.com"})
		if r.calls != 1 || !errors.Is(err, original) {
			t.Fatalf("repository error lost: %v", err)
		}
	}
}
