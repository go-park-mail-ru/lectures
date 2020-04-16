package gqlgen6

//go:generate go run github.com/99designs/gqlgen -v

import (
	"context"
	"crypto/md5"
	"fmt"
	"io/ioutil"
	"log"
	"strconv"
	"time"

	"github.com/99designs/gqlgen/graphql"
)

type Resolver struct {
	PhotosData map[string]*Photo
	Users      map[uint]*User
}

func (r *Resolver) Mutation() MutationResolver {
	return &mutationResolver{r}
}
func (r *Resolver) Photo() PhotoResolver {
	return &photoResolver{r}
}
func (r *Resolver) User() UserResolver {
	return &userResolver{r}
}
func (r *Resolver) Query() QueryResolver {
	return &queryResolver{r}
}

type mutationResolver struct{ *Resolver }

func (r *mutationResolver) RatePhoto(ctx context.Context, id string, direction string) (*Photo, error) {
	log.Println("call mutationResolver.RatePhoto method with id", id, direction)
	rate := 1
	if direction != "up" {
		rate = -1
	}
	ph, ok := r.PhotosData[id]
	if !ok {
		return nil, fmt.Errorf("no photo %v found", id)
	}
	ph.Rating += rate
	return ph, nil
}

func (r *mutationResolver) UploadPhoto(ctx context.Context, comment string, file graphql.Upload) (*Photo, error) {
	sessionUserID := ctx.Value("userID").(uint)

	content, err := ioutil.ReadAll(file.File)
	if err != nil {
		return nil, err
	}
	hasher := md5.New()
	hasher.Write(content)

	log.Printf("incoming file %v, %v bytes, md5 %x", file.Filename, file.Size, hasher.Sum(nil))
	ph := &Photo{
		ID:      42,
		UserID:  sessionUserID,
		Comment: comment,
		URL:     "/photos/" + file.Filename,
	}
	r.PhotosData[strconv.Itoa(int(ph.ID))] = ph
	return ph, nil
}

type userResolver struct{ *Resolver }

func (r *userResolver) Photos(ctx context.Context, obj *User, count int) ([]*Photo, error) {
	log.Println("call userResolver.Photos with count", count)
	id, _ := strconv.Atoi(obj.ID)
	items := []*Photo{}
	for _, ph := range r.PhotosData {
		if ph.UserID != uint(id) {
			continue
		}
		items = append(items, ph)
	}
	return items, nil
}

type photoResolver struct{ *Resolver }

func (r *photoResolver) ID(ctx context.Context, obj *Photo) (string, error) {
	return obj.Id(), nil
}

func (r *photoResolver) User(ctx context.Context, obj *Photo) (*User, error) {
	// return r.Users[obj.UserID], nil
	log.Println("call photoResolver.User", obj.UserID)
	start := time.Now()
	user, err := ctx.Value("userLoaderKey").(*UserLoader).Load(obj.UserID)
	log.Println("get photoResolver.User", obj.UserID, "from UserLoader, time ", time.Since(start))
	return user, err
}

type queryResolver struct{ *Resolver }

func (r *queryResolver) Timeline(ctx context.Context) ([]*Photo, error) {
	log.Println("call queryResolver.Timeline with ctx.userID", ctx.Value("userID"))
	items := []*Photo{}
	for _, ph := range r.PhotosData {
		items = append(items, ph)
	}
	return items, nil
}

func (r *queryResolver) User(ctx context.Context, userID string) (*User, error) {
	log.Println("call queryResolver.User for", userID)
	id, _ := strconv.Atoi(userID)
	return r.Users[uint(id)], nil
}

func (r *queryResolver) Photos(ctx context.Context, userID string) ([]*Photo, error) {
	log.Println("call queryResolver.Photos")
	id, _ := strconv.Atoi(userID)
	items := []*Photo{}
	for _, ph := range r.PhotosData {
		if ph.UserID != uint(id) {
			continue
		}
		items = append(items, ph)
	}
	return items, nil
}
