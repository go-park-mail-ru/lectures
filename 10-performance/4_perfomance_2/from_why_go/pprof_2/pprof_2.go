package main

import (
	"fmt"
	"net/http"
	_ "net/http/pprof"
	"time"
)

type Post struct {
	ID       int
	Text     string
	Author   string
	Comments int
	Time     time.Time
}

func getPost(out chan []Post) {
	posts := []Post{}
	for i := 1; i < 10; i++ {
		post := Post{ID: 1, Text: "text"}
		posts = append(posts, post)
	}
	out <- posts
}

func longHeavyWork(ch chan bool) {
	time.Sleep(1 * time.Minute)
	ch <- true
}

func handleLeak(w http.ResponseWriter, req *http.Request) {
	res := make(chan bool)
	go longHeavyWork(res)
	 <-res
}

func main() {
	http.HandleFunc("/", handleLeak)

	fmt.Println("starting server at :8080")
	http.ListenAndServe(":8080", nil)
}

/*

go build -o pprof_2.exe pprof_2.go && ./pprof_2.exe

hey -z 30s -q 100 http://127.0.0.1:8080/

./pprof_2.sh

*/
