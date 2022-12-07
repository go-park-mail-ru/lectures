package main

import (
	"fmt"
	"log"
	"net"

	"github.com/go-park-mail-ru/lectures/8-microservices/5_grpc_2/session"

	"google.golang.org/grpc"
)

func main() {
	lis, err := net.Listen("tcp", ":8081")
	if err != nil {
		log.Fatalln("cant listet port", err)
	}

	server := grpc.NewServer()

	session.RegisterAuthCheckerServer(server, NewSessionManager())

	fmt.Println("starting server at :8081")
	server.Serve(lis)
}
