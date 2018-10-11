package main

import (
	"log"
	"net/http"

	graphql "github.com/graph-gophers/graphql-go"
	"github.com/graph-gophers/graphql-go/relay"
)

type query struct{}

type User struct {
	ID   int
	Name string
}

func (_ *query) Hello() string { return "Hello, world!" }

func (_ *query) User() User {
	return User{
		ID:   1,
		Name: "Dmitry",
	}
}

func main() {
	s := `
		schema {
			query: Query
		}
		type User {
			ID: Int!
			name: String!
		}
		type Query {
			user: User!
			hello: String!
		}
	`
	schema := graphql.MustParseSchema(s, &query{})
	http.Handle("/query", &relay.Handler{Schema: schema})
	log.Fatal(http.ListenAndServe(":9090", nil))
}
