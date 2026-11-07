package article_test

import (
	"context"
	"errors"
	"lectures/cleanarchreview/internal/entity"
	"lectures/cleanarchreview/internal/usecase/article"
	"testing"
	"time"
)

type authors struct{ calls int }

func (a *authors) GetByID(context.Context, int64) (entity.Author, error) {
	a.calls++
	return entity.NewAuthor(1, "Автор")
}

type repo struct {
	article.Repository
	saved bool
	err   error
}

func (r *repo) Save(context.Context, entity.Article) error { r.saved = true; return r.err }

type clock struct{}

func (clock) Now() time.Time { return time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC) }

type ids struct{}

func (ids) Next() int64 { return 1 }

func TestCreateValidatesBeforeSave(t *testing.T) {
	r := &repo{}
	s := article.New(r, &authors{}, clock{}, ids{})
	_, err := s.Create(context.Background(), article.CreateInput{AuthorID: 1, Title: " ", Content: "Текст"})
	if !errors.Is(err, entity.ErrInvalidArticle) || r.saved {
		t.Fatalf("invalid entity saved: %v", err)
	}
}

func TestCreatePreservesStorageConflict(t *testing.T) {
	r := &repo{err: article.ErrConflict}
	s := article.New(r, &authors{}, clock{}, ids{})
	_, err := s.Create(context.Background(), article.CreateInput{AuthorID: 1, Title: "Статья", Content: "Текст"})
	if !r.saved || !errors.Is(err, article.ErrConflict) {
		t.Fatalf("conflict lost: %v", err)
	}
}
