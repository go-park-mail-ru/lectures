package main

import (
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"

	"github.com/go-park-mail-ru/lectures/8-microservices/5_grpc/session"
)

func main() {
	listener, err := net.Listen("tcp", ":8081")
	if err != nil {
		log.Fatalln("can't listen port", err)
	}

	server := grpc.NewServer()
	session.RegisterAuthCheckerServer(server, NewSessionManager())

	fmt.Println("starting server at :8081")
	log.Fatal(server.Serve(listener))
}
