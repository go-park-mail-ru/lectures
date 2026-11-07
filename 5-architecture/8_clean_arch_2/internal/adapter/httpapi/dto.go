package httpapi

import (
	"lectures/cleanarchreview/internal/usecase/article"
	"time"
)

type createRequest struct {
	AuthorID int64  `json:"author_id"`
	Title    string `json:"title"`
	Content  string `json:"content"`
}
type updateRequest struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}
type authorResponse struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}
type articleResponse struct {
	ID        int64          `json:"id"`
	Title     string         `json:"title"`
	Content   string         `json:"content"`
	Author    authorResponse `json:"author"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
}

func toResponse(v article.View) articleResponse {
	a := v.Article
	return articleResponse{ID: a.ID(), Title: a.Title(), Content: a.Content(),
		Author:    authorResponse{ID: v.Author.ID(), Name: v.Author.Name()},
		CreatedAt: a.CreatedAt(), UpdatedAt: a.UpdatedAt()}
}
