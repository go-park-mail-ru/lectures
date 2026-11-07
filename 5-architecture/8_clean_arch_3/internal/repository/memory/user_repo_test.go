package memory_test

import (
	"context"
	"errors"
	"testing"

	"lectures/cleanarch3/internal/domain"
	"lectures/cleanarch3/internal/repository/memory"
	"lectures/cleanarch3/internal/usecase"
)

func TestConcurrentRegistrationKeepsEmailUnique(t *testing.T) {
	r := memory.NewUserRepository()
	u, err := domain.NewUser("Анна", "anna@example.com")
	if err != nil {
		t.Fatal(err)
	}
	const attempts = 20
	results := make(chan error, attempts)
	start := make(chan struct{})
	for i := 0; i < attempts; i++ {
		go func() {
			<-start
			_, err := r.Create(context.Background(), u)
			results <- err
		}()
	}
	close(start)
	successes := 0
	for i := 0; i < attempts; i++ {
		err := <-results
		if err == nil {
			successes++
		} else if !errors.Is(err, usecase.ErrEmailTaken) {
			t.Fatalf("unexpected registration error: %v", err)
		}
	}
	items, err := r.List(context.Background(), 0, attempts)
	if err != nil || successes != 1 || len(items) != 1 {
		t.Fatalf("duplicate email: successes=%d, users=%d, err=%v", successes, len(items), err)
	}
}

func TestProfileConflictAndEmailReuse(t *testing.T) {
	ctx := context.Background()
	r := memory.NewUserRepository()
	uc := usecase.NewUserUseCase(r)
	anna, err := uc.Register(ctx, usecase.CreateUserInput{Name: "Анна", Email: "anna@example.com"})
	if err != nil {
		t.Fatal(err)
	}
	boris, err := uc.Register(ctx, usecase.CreateUserInput{Name: "Борис", Email: "boris@example.com"})
	if err != nil {
		t.Fatal(err)
	}
	_, err = uc.Update(ctx, usecase.UpdateUserInput{ID: anna.ID(), Name: "Анна", Email: boris.Email()})
	if !errors.Is(err, usecase.ErrEmailTaken) {
		t.Fatalf("duplicate email accepted in update: %v", err)
	}
	stored, err := uc.GetByID(ctx, anna.ID())
	if err != nil || stored.Email() != anna.Email() {
		t.Fatalf("conflict changed stored user: %v, %v", stored, err)
	}
	_, err = uc.Update(ctx, usecase.UpdateUserInput{ID: anna.ID(), Name: "Анна", Email: "new@example.com"})
	if err != nil {
		t.Fatal(err)
	}
	if _, err := uc.Register(ctx, usecase.CreateUserInput{Name: "Другой пользователь", Email: anna.Email()}); err != nil {
		t.Fatalf("old email remains occupied: %v", err)
	}
	if err := uc.Delete(ctx, boris.ID()); err != nil {
		t.Fatal(err)
	}
	if _, err := uc.Register(ctx, usecase.CreateUserInput{Name: "Другой Борис", Email: boris.Email()}); err != nil {
		t.Fatalf("deleted user's email remains occupied: %v", err)
	}
}

func TestCanceledRegistrationDoesNotSave(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	r := memory.NewUserRepository()
	uc := usecase.NewUserUseCase(r)
	_, err := uc.Register(ctx, usecase.CreateUserInput{Name: "Анна", Email: "anna@example.com"})
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("cancellation ignored: %v", err)
	}
	items, err := r.List(context.Background(), 0, 10)
	if err != nil || len(items) != 0 {
		t.Fatalf("canceled registration saved user: %v, %v", items, err)
	}
}
