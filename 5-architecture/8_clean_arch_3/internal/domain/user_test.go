package domain_test

import (
	"errors"
	"testing"

	"lectures/cleanarch3/internal/domain"
)

func TestUserRules(t *testing.T) {
	u, err := domain.NewUser("  Анна  ", "  Anna@Example.com  ")
	if err != nil || u.ID() != 0 || u.Name() != "Анна" || u.Email() != "anna@example.com" {
		t.Fatalf("new user: %v, %v", u, err)
	}
	for _, in := range []struct{ name, email string }{
		{" ", "anna@example.com"},
		{"Анна", ""},
		{"Анна", "not-an-email"},
		{"Анна", "Anna <anna@example.com>"},
	} {
		if _, err := domain.NewUser(in.name, in.email); !errors.Is(err, domain.ErrInvalidUser) {
			t.Fatalf("invalid profile accepted: %+v, %v", in, err)
		}
	}
	if _, err := domain.RestoreUser(0, "Анна", "anna@example.com"); !errors.Is(err, domain.ErrInvalidUser) {
		t.Fatalf("stored user without ID accepted: %v", err)
	}
}

func TestChangeProfilePreservesIdentityAndOriginal(t *testing.T) {
	u, err := domain.RestoreUser(42, "Анна", "anna@example.com")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := u.ChangeProfile(" ", "other@example.com"); !errors.Is(err, domain.ErrInvalidUser) {
		t.Fatalf("invalid change accepted: %v", err)
	}
	changed, err := u.ChangeProfile("Анна Иванова", "other@example.com")
	if err != nil || changed.ID() != 42 || changed.Name() != "Анна Иванова" {
		t.Fatalf("profile change: %v, %v", changed, err)
	}
	if u.Name() != "Анна" || u.Email() != "anna@example.com" {
		t.Fatalf("original profile changed: %v", u)
	}
}
