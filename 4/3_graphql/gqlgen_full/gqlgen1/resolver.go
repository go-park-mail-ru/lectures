package gqlgen1

import (
	"context"
) // THIS CODE IS A STARTING POINT ONLY. IT WILL NOT BE UPDATED WITH SCHEMA CHANGES.

type Resolver struct{}

func (r *Resolver) Mutation() MutationResolver {
	return &mutationResolver{r}
}
func (r *Resolver) Query() QueryResolver {
	return &queryResolver{r}
}

type mutationResolver struct{ *Resolver }

func (r *mutationResolver) RatePhoto(ctx context.Context, photoID string, direction string) (*Photo, error) {
	panic("not implemented")
}

type queryResolver struct{ *Resolver }

func (r *queryResolver) Timeline(ctx context.Context) ([]*Photo, error) {
	panic("not implemented")
}
func (r *queryResolver) User(ctx context.Context, userID string) (*User, error) {
	panic("not implemented")
}
func (r *queryResolver) Photos(ctx context.Context, userID string) ([]*Photo, error) {
	panic("not implemented")
}
