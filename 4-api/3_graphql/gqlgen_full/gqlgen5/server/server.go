package main

import (
	"context"
	gqlgen "gqlgen5"
	"log"
	"net/http"
	"time"

	"github.com/99designs/gqlgen/handler"
)

/*
curl localhost:8080/query \
  -F operations='{ "query": "mutation($comment: String!, $file: Upload!) { uploadPhoto(comment: $comment, file: $file) { id } }", "variables": { "comment": "building 5 comment", "file": null } }' \
  -F map='{ "0": ["variables.file"] }' \
  -F 0=@./test_file.txt \
  --trace-ascii -


  {
  query: `
    mutation($comment: String!, $file: Upload!) {
      uploadPhoto(comment: $comment, file: $file) {
        id
      }
    }
  `,
  variables: {
	comment: "building 5 comment",
    file: File // test_file.txt
  }
}



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
		ID:      1,
		UserID:  1,
		URL:     "https://via.placeholder.com/300",
		Comment: "from studio",
		Rating:  1,
		Liked:   true,
	},
	// "2": {
	// 	ID:      2,
	// 	UserID:  1,
	// 	URL:     "https://via.placeholder.com/300",
	// 	Comment: "cool view",
	// 	Rating:  0,
	// 	Liked:   false,
	// },
	"3": {
		ID:      3,
		UserID:  2,
		URL:     "https://via.placeholder.com/300",
		Comment: "at work",
		Rating:  0,
		Liked:   false,
	},
}

// go run github.com/vektah/dataloaden UserLoader uint *coursera/3p/graphql/gqlgen3.User

func UserLoaderMiddleware(resolver *gqlgen.Resolver, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg := gqlgen.UserLoaderConfig{
			MaxBatch: 100,
			Wait:     1 * time.Millisecond,
			Fetch: func(ids []uint) ([]*gqlgen.User, []error) {
				// имеем доступ до r *http.Request - там context с сессией пользователя
				sessionUserID := r.Context().Value("userID").(uint)
				log.Printf("UserLoader Request - ids %v for user %v\n", ids, sessionUserID)

				log.Printf("request %v\n", r)
				log.Printf("ctx %v\n", r.Context())

				users := make([]*gqlgen.User, len(ids))
				for i, id := range ids {
					// имеем доступ до resolver *gqlgen.Resolver - там коннет к базе
					users[i] = resolver.Users[id]
				}
				return users, nil
			},
		}
		userLoader := gqlgen.NewUserLoader(cfg)
		ctx := context.WithValue(r.Context(), "userLoaderKey", userLoader)
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// log.Println("new request")
		ctx := context.WithValue(r.Context(), "userID", uint(1))
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}

func main() {
	http.Handle("/", handler.Playground("GraphQL playground", "/query"))

	resolver := &gqlgen.Resolver{
		Users:      users,
		PhotosData: photos,
	}
	cfg := gqlgen.Config{
		Resolvers: resolver,
	}
	cfg.Complexity.User.Photos = func(childComplexity, count int) int {
		return count * childComplexity
	}

	gqlHandler := handler.GraphQL(
		gqlgen.NewExecutableSchema(cfg),
		handler.ComplexityLimit(500),
	)
	handler := UserLoaderMiddleware(resolver, gqlHandler)
	handler = AuthMiddleware(handler)
	http.Handle("/query", handler)

	port := "8080"
	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
