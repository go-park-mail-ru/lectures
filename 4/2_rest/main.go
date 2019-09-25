package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type Result struct {
	Body interface{} `json:"body,omitempty"`
	Err  string      `json:"err,omitempty"`
}

type BooksHandler struct {
	store *BookStore
}

func (api *BooksHandler) List(w http.ResponseWriter, r *http.Request) {

	books, err := api.store.GetBooks()
	if err != nil {
		http.Error(w, `{"error":"db"}`, 500)
		return
	}

	body := map[string]interface{}{
		"books": books,
	}
	json.NewEncoder(w).Encode(&Result{Body: body})
}

// http://127.0.0.1:8080/add?title=test&price=123

func (api *BooksHandler) Add(w http.ResponseWriter, r *http.Request) {
	title := r.FormValue("title")
	price, _ := strconv.Atoi(r.FormValue("price"))

	in := &Book{
		Title: title,
		Price: uint(price),
	}
	id, err := api.store.AddBook(in)
	if err != nil {
		http.Error(w, `{"error":"db"}`, 500)
		return
	}

	body := map[string]interface{}{
		"id": id,
	}
	json.NewEncoder(w).Encode(&Result{Body: body})
}

func main() {
	r := mux.NewRouter()

	api := &BooksHandler{
		store: NewBookStore(),
	}

	r.HandleFunc("/", api.List)
	r.HandleFunc("/add", api.Add)

	log.Println("start serving :8080")
	http.ListenAndServe(":8080", r)
}
