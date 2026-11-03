package article

import (
	"context"
	"fmt"
	"lectures/cleanarchreview/internal/entity"
	"strings"
)

type Service struct {
	articles Repository
	authors  Authors
	clock    Clock
	ids      IDs
}

func New(articles Repository, authors Authors, clock Clock, ids IDs) *Service {
	return &Service{articles: articles, authors: authors, clock: clock, ids: ids}
}

func (s *Service) Create(ctx context.Context, in CreateInput) (View, error) {
	author, err := s.authors.GetByID(ctx, in.AuthorID)
	if err != nil {
		return View{}, fmt.Errorf("find author: %w", err)
	}
	a, err := entity.NewArticle(s.ids.Next(), author.ID(), in.Title, in.Content, s.clock.Now())
	if err != nil {
		return View{}, err
	}

	if err := s.articles.Save(ctx, a); err != nil {
		return View{}, fmt.Errorf("save article: %w", err)
	}
	return View{Article: a, Author: author}, nil
}

func (s *Service) GetByID(ctx context.Context, id int64) (View, error) {
	if id <= 0 {
		return View{}, ErrInvalidInput
	}
	a, err := s.articles.GetByID(ctx, id)
	if err != nil {
		return View{}, err
	}
	return s.view(ctx, a)
}

func (s *Service) GetByTitle(ctx context.Context, title string) (View, error) {
	title = strings.TrimSpace(title)
	if title == "" {
		return View{}, ErrInvalidInput
	}
	a, err := s.articles.GetByTitle(ctx, title)
	if err != nil {
		return View{}, err
	}
	return s.view(ctx, a)
}

func (s *Service) Update(ctx context.Context, in UpdateInput) (View, error) {
	if in.ID <= 0 {
		return View{}, ErrInvalidInput
	}
	a, err := s.articles.GetByID(ctx, in.ID)
	if err != nil {
		return View{}, err
	}
	a, err = a.Revise(in.Title, in.Content, s.clock.Now())
	if err != nil {
		return View{}, err
	}
	view, err := s.view(ctx, a)
	if err != nil {
		return View{}, err
	}
	if err := s.articles.Save(ctx, a); err != nil {
		return View{}, err
	}
	return view, nil
}

func (s *Service) Fetch(ctx context.Context, afterID int64, limit int) (Page, error) {
	if afterID < 0 || limit < 0 || limit > 100 {
		return Page{}, ErrInvalidInput
	}
	if limit == 0 {
		limit = 10
	}

	items, err := s.articles.List(ctx, afterID, limit+1)
	if err != nil {
		return Page{}, err
	}
	page := Page{Items: make([]View, 0, limit)}
	if len(items) > limit {
		items = items[:limit]
		page.NextID = items[limit-1].ID()
	}
	for _, a := range items {
		v, err := s.view(ctx, a)
		if err != nil {
			return Page{}, err
		}
		page.Items = append(page.Items, v)
	}
	return page, nil
}

func (s *Service) Delete(ctx context.Context, id int64) error {
	if id <= 0 {
		return ErrInvalidInput
	}
	return s.articles.Delete(ctx, id)
}

func (s *Service) view(ctx context.Context, a entity.Article) (View, error) {
	author, err := s.authors.GetByID(ctx, a.AuthorID())
	if err != nil {
		return View{}, err
	}
	return View{Article: a, Author: author}, nil
}
