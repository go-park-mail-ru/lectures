package entity

import (
	"errors"
	"testing"
	"time"
)

func TestArticleInvariants(t *testing.T) {
	now := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	a, err := NewArticle(1, 2, "  Заголовок  ", "Текст", now)
	if err != nil || a.Title() != "Заголовок" {
		t.Fatalf("new article: %v, %v", a, err)
	}
	if _, err := a.Revise(" ", "Текст", now); !errors.Is(err, ErrInvalidArticle) {
		t.Fatalf("empty title accepted: %v", err)
	}
	if a.Title() != "Заголовок" {
		t.Fatal("failed revision changed original")
	}
	updated, err := a.Revise("Новый", "Текст", now.Add(time.Hour))
	if err != nil || updated.CreatedAt() != now || updated.ID() != a.ID() {
		t.Fatalf("revision lost identity: %v", err)
	}
	if _, err := updated.Revise("Ещё", "Текст", now); !errors.Is(err, ErrInvalidArticle) {
		t.Fatal("time moved backwards")
	}
}
