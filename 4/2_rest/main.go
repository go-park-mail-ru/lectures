package main

import (
	"net/http"

	"github.com/gorilla/mux"
)

type Author struct {
	Name string `json:"name"`
}

type Book struct {
	ID     uint   `json:"id"`
	Title  string `json:"string"`
	Author uint   `json:"author"`
}

var books []Book

func RootHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("test"))
}

func main() {
	r := mux.NewRouter()

	r.HandleFunc("/", RootHandler)

	http.ListenAndServe(":8080", r)
}
