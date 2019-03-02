package main

import (
	"net/http"

	"github.com/gorilla/mux"
)

type User struct {
	ID       uint   `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
}

var users []User

func RootHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("test"))
}

func main() {
	r := mux.NewRouter()

	r.HandleFunc("/", RootHandler)

	http.ListenAndServe(":8080", r)
}
