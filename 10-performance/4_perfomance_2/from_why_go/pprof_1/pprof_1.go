package main

import (
	"bytes"
	"fmt"
	"net/http"
	_ "net/http/pprof"
	"sync"
	"time"
)

type Post struct {
	ID       int
	Text     string
	Author   string
	Comments int
	Time     time.Time
}

func handleSlow(w http.ResponseWriter, req *http.Request) {
	s := ""
	for i := 0; i < 1000; i++ {
		p := &Post{ID: i, Text: "new post"}
		s += fmt.Sprintf("%#v", p)
	}
	w.Write([]byte(s))
}

func main() {
	http.HandleFunc("/", handleSlow)
	http.HandleFunc("/fast", handleFast)

	fmt.Println("starting server at :8080")
	fmt.Println(http.ListenAndServe(":8080", nil))
}

/*
go build -o pprof_1.exe pprof_1.go && ./pprof_1.exe

hey -z 20s http://127.0.0.1:8080
hey -z 20s http://127.0.0.1:8080/fast

curl http://127.0.0.1:8080/debug/pprof/heap -o mem_out.txt
curl http://127.0.0.1:8080/debug/pprof/profile?seconds=5 -o cpu_out.txt

go tool pprof -svg -inuse_space pprof_1.exe mem_out.txt > mem_is.svg
go tool pprof -svg -inuse_objects pprof_1.exe mem_out.txt > mem_oo.svg
go tool pprof -svg -alloc_space pprof_1.exe mem_out.txt > mem_as.svg
go tool pprof -svg -alloc_objects pprof_1.exe mem_out.txt > mem_ao.svg
go tool pprof -svg pprof_1.exe cpu_out.txt > cpu.svg

*/

var dataPool = sync.Pool{
	New: func() interface{} {
		return &Post{}
	},
}

var bufPool = sync.Pool{
	New: func() interface{} {
		return bytes.NewBuffer(make([]byte, 0, 50000))
	},
}

func handleFast(w http.ResponseWriter, req *http.Request) {
	buf := bufPool.Get().(*bytes.Buffer)
	for i := 0; i < 1000; i++ {
		p := dataPool.Get().(*Post)
		p.ID = i
		p.Text = "new post"
		fmt.Fprintf(buf, "%#v", p)
		*p = Post{}
		dataPool.Put(p)
	}
	buf.Reset()
	bufPool.Put(buf)
	w.Write(buf.Bytes())
}
