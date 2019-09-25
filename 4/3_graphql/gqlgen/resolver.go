package gqlgen

import (
	"context"
) // THIS CODE IS A STARTING POINT ONLY. IT WILL NOT BE UPDATED WITH SCHEMA CHANGES.

type Resolver struct{}

func (r *Resolver) Query() QueryResolver {
	return &queryResolver{r}
}

type queryResolver struct{ *Resolver }

func (r *queryResolver) Books(ctx context.Context) ([]*Book, error) {
	return []*Book{
		{
			ID:     "1",
			Title:  "The Moon is a harsh mistress",
			Price:  200,
			Author: &Author{Name: "John"},
		},
	}, nil
}
