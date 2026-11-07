// Package entity содержит бизнес-объекты, а не модели JSON или таблиц.
package entity

import (
	"errors"
	"strings"
	"time"
)

var ErrInvalidArticle = errors.New("invalid article")

type Article struct {
	id, authorID         int64
	title, content       string
	createdAt, updatedAt time.Time
}

func NewArticle(id, authorID int64, title, content string, now time.Time) (Article, error) {
	return RestoreArticle(id, authorID, title, content, now, now)
}

func RestoreArticle(id, authorID int64, title, content string, created, updated time.Time) (Article, error) {
	title, content = strings.TrimSpace(title), strings.TrimSpace(content)
	if id <= 0 || authorID <= 0 || title == "" || content == "" || created.IsZero() || updated.Before(created) {
		return Article{}, ErrInvalidArticle
	}
	return Article{id: id, authorID: authorID, title: title, content: content, createdAt: created, updatedAt: updated}, nil
}

func (a Article) Revise(title, content string, now time.Time) (Article, error) {
	if now.Before(a.updatedAt) {
		return Article{}, ErrInvalidArticle
	}
	return RestoreArticle(a.id, a.authorID, title, content, a.createdAt, now)
}

func (a Article) ID() int64            { return a.id }
func (a Article) AuthorID() int64      { return a.authorID }
func (a Article) Title() string        { return a.title }
func (a Article) Content() string      { return a.content }
func (a Article) CreatedAt() time.Time { return a.createdAt }
func (a Article) UpdatedAt() time.Time { return a.updatedAt }
