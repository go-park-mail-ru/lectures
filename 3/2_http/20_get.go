package main

import (
	"fmt"
	"net/http"
)

func handler(w http.ResponseWriter, r *http.Request) {
	for key, value := range r.URL.Query() {
		fmt.Println(key, value)
	}
}

func main() {
	http.HandleFunc("/", handler)

	fmt.Println("starting server at :8080")
	http.ListenAndServe(":8080", nil)
}
