package gql

import (
	"context"
)

type Resolver struct{}

func (r *Resolver) Query() QueryResolver {
	return &queryResolver{r}
}

type queryResolver struct{ *Resolver }

func (r *queryResolver) Books(ctx context.Context) ([]Book, error) {
	return []Book{
		Book{
			Title: "The Moon is a harsh mistress",
			Price: 200,
		},
	}, nil
}
