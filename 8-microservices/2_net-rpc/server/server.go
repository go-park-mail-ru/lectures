package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"net/rpc"
)

func main() {
	sessionManager := NewSessionManager()

	rpc.Register(sessionManager)
	rpc.HandleHTTP()

	listener, e := net.Listen("tcp", ":8081")
	if e != nil {
		log.Fatal("listen error:", e)
	}

	fmt.Println("starting server at :8081")
	log.Fatal(http.Serve(listener, nil))
}
