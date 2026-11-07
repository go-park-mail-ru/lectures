// Package mysql показывает границу SQL. Подключение и драйвер выбирает main.
// GetByID реализован как образец маппинга; остальные методы — явные заглушки.
package mysql

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"lectures/cleanarchreview/internal/entity"
	"lectures/cleanarchreview/internal/usecase/article"
	"time"
)

// articleRow — модель хранения. Типы sql.* не выходят из адаптера.
type articleRow struct {
	id, authorID     int64
	title, content   string
	created, updated time.Time
}

func (r articleRow) toEntity() (entity.Article, error) {
	return entity.RestoreArticle(r.id, r.authorID, r.title, r.content, r.created, r.updated)
}

type Articles struct{ db *sql.DB }

func NewArticles(db *sql.DB) *Articles { return &Articles{db: db} }

var _ article.Repository = (*Articles)(nil)

func (r *Articles) GetByID(ctx context.Context, id int64) (entity.Article, error) {
	var row articleRow
	err := r.db.QueryRowContext(ctx, `SELECT id, author_id, title, content, created_at, updated_at FROM article WHERE id = ?`, id).
		Scan(&row.id, &row.authorID, &row.title, &row.content, &row.created, &row.updated)
	if errors.Is(err, sql.ErrNoRows) {
		return entity.Article{}, article.ErrNotFound
	}
	if err != nil {
		return entity.Article{}, fmt.Errorf("select article: %w", err)
	}
	return row.toEntity()
}

func (*Articles) Save(context.Context, entity.Article) error { return article.ErrNotImplemented }
func (*Articles) GetByTitle(context.Context, string) (entity.Article, error) {
	return entity.Article{}, article.ErrNotImplemented
}
func (*Articles) List(context.Context, int64, int) ([]entity.Article, error) {
	return nil, article.ErrNotImplemented
}
func (*Articles) Delete(context.Context, int64) error { return article.ErrNotImplemented }

type Authors struct{ db *sql.DB }

func NewAuthors(db *sql.DB) *Authors { return &Authors{db: db} }

var _ article.Authors = (*Authors)(nil)

func (r *Authors) GetByID(ctx context.Context, id int64) (entity.Author, error) {
	var name string
	err := r.db.QueryRowContext(ctx, `SELECT name FROM author WHERE id = ?`, id).Scan(&name)
	if errors.Is(err, sql.ErrNoRows) {
		return entity.Author{}, article.ErrNotFound
	}
	if err != nil {
		return entity.Author{}, fmt.Errorf("select author: %w", err)
	}
	return entity.NewAuthor(id, name)
}
