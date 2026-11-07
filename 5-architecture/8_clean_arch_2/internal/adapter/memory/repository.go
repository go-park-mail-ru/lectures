// Package memory — выходной адаптер. Карты являются деталью хранения.
package memory

import (
	"context"
	"lectures/cleanarchreview/internal/entity"
	"lectures/cleanarchreview/internal/usecase/article"
	"sort"
	"sync"
)

type Articles struct {
	mu    sync.RWMutex
	items map[int64]entity.Article
}

func NewArticles() *Articles { return &Articles{items: make(map[int64]entity.Article)} }

var _ article.Repository = (*Articles)(nil)

func (r *Articles) Save(ctx context.Context, a entity.Article) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	// Нулевое значение Go тоже возможно: защищаем границу хранения.
	if _, err := entity.RestoreArticle(a.ID(), a.AuthorID(), a.Title(), a.Content(), a.CreatedAt(), a.UpdatedAt()); err != nil {
		return err
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	for id, other := range r.items {
		if id != a.ID() && other.Title() == a.Title() {
			return article.ErrConflict
		}
	}
	r.items[a.ID()] = a
	return nil
}

func (r *Articles) GetByID(ctx context.Context, id int64) (entity.Article, error) {
	if err := ctx.Err(); err != nil {
		return entity.Article{}, err
	}
	r.mu.RLock()
	defer r.mu.RUnlock()
	a, ok := r.items[id]
	if !ok {
		return entity.Article{}, article.ErrNotFound
	}
	return a, nil
}

func (r *Articles) GetByTitle(ctx context.Context, title string) (entity.Article, error) {
	if err := ctx.Err(); err != nil {
		return entity.Article{}, err
	}
	r.mu.RLock()
	defer r.mu.RUnlock()
	for _, a := range r.items {
		if a.Title() == title {
			return a, nil
		}
	}
	return entity.Article{}, article.ErrNotFound
}

func (r *Articles) List(ctx context.Context, afterID int64, limit int) ([]entity.Article, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	if afterID < 0 || limit <= 0 {
		return nil, article.ErrInvalidInput
	}
	r.mu.RLock()
	defer r.mu.RUnlock()
	items := make([]entity.Article, 0)
	for id, a := range r.items {
		if id > afterID {
			items = append(items, a)
		}
	}
	sort.Slice(items, func(i, j int) bool { return items[i].ID() < items[j].ID() })
	if len(items) > limit {
		items = items[:limit]
	}
	return items, nil
}

func (r *Articles) Delete(ctx context.Context, id int64) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.items[id]; !ok {
		return article.ErrNotFound
	}
	delete(r.items, id)
	return nil
}

// Authors — неизменяемый после создания справочник для демонстрации.
type Authors struct{ items map[int64]entity.Author }

func NewAuthors(authors ...entity.Author) *Authors {
	r := &Authors{items: make(map[int64]entity.Author)}
	for _, a := range authors {
		r.items[a.ID()] = a
	}
	return r
}

var _ article.Authors = (*Authors)(nil)

func (r *Authors) GetByID(ctx context.Context, id int64) (entity.Author, error) {
	if err := ctx.Err(); err != nil {
		return entity.Author{}, err
	}
	a, ok := r.items[id]
	if !ok {
		return entity.Author{}, article.ErrNotFound
	}
	return a, nil
}
