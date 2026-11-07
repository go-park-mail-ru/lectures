// Package article реализует сценарии приложения и владеет выходными портами.
package article

import (
	"context"
	"errors"
	"lectures/cleanarchreview/internal/entity"
	"time"
)

var (
	ErrNotFound       = errors.New("not found")
	ErrConflict       = errors.New("article title already exists")
	ErrInvalidInput   = errors.New("invalid input")
	ErrNotImplemented = errors.New("adapter not implemented")
)

type Repository interface {
	Save(context.Context, entity.Article) error
	GetByID(context.Context, int64) (entity.Article, error)
	GetByTitle(context.Context, string) (entity.Article, error)
	List(context.Context, int64, int) ([]entity.Article, error)
	Delete(context.Context, int64) error
}

type Authors interface {
	GetByID(context.Context, int64) (entity.Author, error)
}

type Clock interface{ Now() time.Time }
type IDs interface{ Next() int64 }

// Эти DTO принадлежат приложению: у них нет JSON-тегов
type CreateInput struct {
	AuthorID       int64
	Title, Content string
}
type UpdateInput struct {
	ID             int64
	Title, Content string
}
type View struct {
	Article entity.Article
	Author  entity.Author
}
type Page struct {
	Items  []View
	NextID int64
}
