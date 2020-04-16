package main

import (
	"context"
	gqlgen "gqlgen2"
	"log"
	"net/http"

	"github.com/99designs/gqlgen/handler"
)

/*
query{timeline{id,url,user{id,name}}}
query{user(userID:"1"){id,avatar, name}}
mutation _{ratePhoto(photoID:"1", direction:"up"){id,url,rating,user{id,name}}}
*/

var users = map[uint]*gqlgen.User{
	1: {
		ID:     "1",
		Name:   "rvasily",
		Avatar: "https://via.placeholder.com/150",
	},
	2: {
		ID:     "2",
		Name:   "v.romanov",
		Avatar: "https://via.placeholder.com/150",
	},
}

var photos = map[string]*gqlgen.Photo{
	"1": {
		ID:       1,
		UserID:   1,
		URL:      "https://via.placeholder.com/300",
		Comment:  "fromn studio",
		Rating:   1,
		Liked:    true,
		Followed: false,
	},
	"2": {
		ID:       2,
		UserID:   1,
		URL:      "https://via.placeholder.com/300",
		Comment:  "cool view",
		Rating:   0,
		Liked:    false,
		Followed: false,
	},
	"3": {
		ID:       3,
		UserID:   2,
		URL:      "https://via.placeholder.com/300",
		Comment:  "at work",
		Rating:   0,
		Liked:    false,
		Followed: false,
	},
}

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("new request")
		ctx := context.WithValue(r.Context(), "userID", 1)
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}

func main() {
	http.Handle("/", handler.Playground("GraphQL playground", "/query"))

	cfg := gqlgen.Config{
		Resolvers: &gqlgen.Resolver{
			Users:      users,
			PhotosData: photos,
		},
	}
	gqlHandler := handler.GraphQL(gqlgen.NewExecutableSchema(cfg))
	handler := AuthMiddleware(gqlHandler)
	http.Handle("/query", handler)

	port := "8080"
	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
