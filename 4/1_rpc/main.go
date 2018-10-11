package main

import (
	"fmt"
	"net/http"
)

func login(w http.ResponseWriter, r *http.Request) {
	fmt.Println("login")
	w.Write([]byte("login"))
}

func signup(w http.ResponseWriter, r *http.Request) {
	fmt.Println("signup")
	w.Write([]byte("signup"))
}

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		method := r.FormValue("method")

		switch method {
		case "login":
			login(w, r)

		case "signup":
			signup(w, r)
		}
	})
	http.ListenAndServe(":9090", nil)
}
