package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"time"

	"github.com/graphql-go/graphql"
)

type Author struct {
	Name string `json:"string"`
}

type Book struct {
	ID     int64   `json:"id"`
	Title  string  `json:"title"`
	Author uint    `json:"author"`
	Price  float64 `json:"price"`
}

var authors = map[uint]Author{
	1: Author{
		Name: "Robert Heinlein",
	},
}

var books = []Book{
	Book{
		ID:     1,
		Title:  "The Moon is a harsh mistress",
		Author: 1,
		Price:  200,
	},
}

var authorType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "Author",
		Fields: graphql.Fields{
			"name": &graphql.Field{
				Type: graphql.String,
			},
		},
	},
)

var bookType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "Book",
		Fields: graphql.Fields{
			"id": &graphql.Field{
				Type: graphql.Int,
			},
			"title": &graphql.Field{
				Type: graphql.String,
			},
			"author": &graphql.Field{
				Type: authorType,
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
					book, ok := params.Source.(Book)
					if !ok {
						return nil, fmt.Errorf("cannot convert source to Book")
					}

					return authors[book.Author], nil
				},
			},
			"price": &graphql.Field{
				Type: graphql.Float,
			},
		},
	},
)

var mutationType = graphql.NewObject(graphql.ObjectConfig{
	Name: "Mutation",
	Fields: graphql.Fields{
		"create": &graphql.Field{
			Type:        bookType,
			Description: "Create new book",
			Args: graphql.FieldConfigArgument{
				"title": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.String),
				},
				"author": &graphql.ArgumentConfig{
					Type: authorType,
				},
				"price": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.Float),
				},
			},
			Resolve: func(params graphql.ResolveParams) (interface{}, error) {
				rand.Seed(time.Now().UnixNano())

				book := Book{
					ID:     int64(rand.Intn(100000)), // генерируем случайный ID
					Title:  params.Args["title"].(string),
					Author: params.Args["author"].(uint),
					Price:  params.Args["price"].(float64),
				}
				books = append(books, book)
				return book, nil
			},
		},

		"update": &graphql.Field{
			Type:        bookType,
			Description: "Update book by id",
			Args: graphql.FieldConfigArgument{
				"id": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.Int),
				},
				"title": &graphql.ArgumentConfig{
					Type: graphql.String,
				},
				"author": &graphql.ArgumentConfig{
					Type: authorType,
				},
				"price": &graphql.ArgumentConfig{
					Type: graphql.Float,
				},
			},
			Resolve: func(params graphql.ResolveParams) (interface{}, error) {
				id, _ := params.Args["id"].(int)
				title, titleOk := params.Args["title"].(string)
				author, authorOk := params.Args["author"].(uint)
				price, priceOk := params.Args["price"].(float64)
				book := Book{}
				for i, p := range books {
					if int64(id) == p.ID {
						if titleOk {
							books[i].Title = title
						}
						if authorOk {
							books[i].Author = author
						}
						if priceOk {
							books[i].Price = price
						}
						book = books[i]
						break
					}
				}
				return book, nil
			},
		},

		"delete": &graphql.Field{
			Type:        bookType,
			Description: "Delete book by id",
			Args: graphql.FieldConfigArgument{
				"id": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.Int),
				},
			},
			Resolve: func(params graphql.ResolveParams) (interface{}, error) {
				id, _ := params.Args["id"].(int)
				book := Book{}
				for i, p := range books {
					if int64(id) == p.ID {
						book = books[i]
						books = append(books[:i], books[i+1:]...)
					}
				}

				return book, nil
			},
		},
	},
})

var queryType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "Query",
		Fields: graphql.Fields{
			"book": &graphql.Field{
				Type:        bookType,
				Description: "Get book by id",
				Args: graphql.FieldConfigArgument{
					"id": &graphql.ArgumentConfig{
						Type: graphql.Int,
					},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					id, ok := p.Args["id"].(int)
					if ok {
						for _, book := range books {
							if int(book.ID) == id {
								return book, nil
							}
						}
					}
					return nil, nil
				},
			},
			"books": &graphql.Field{
				Type:        graphql.NewList(bookType),
				Description: "Get books list",
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
					return books, nil
				},
			},
		},
	})

var schema, _ = graphql.NewSchema(
	graphql.SchemaConfig{
		Query:    queryType,
		Mutation: mutationType,
	},
)

func executeQuery(query string, schema graphql.Schema) *graphql.Result {
	result := graphql.Do(graphql.Params{
		Schema:        schema,
		RequestString: query,
	})
	if len(result.Errors) > 0 {
		fmt.Printf("errors: %v", result.Errors)
	}
	return result
}

func main() {
	http.HandleFunc("/graphql", func(w http.ResponseWriter, r *http.Request) {
		result := executeQuery(r.URL.Query().Get("query"), schema)
		json.NewEncoder(w).Encode(result)
	})
	http.ListenAndServe(":8080", nil)
}
