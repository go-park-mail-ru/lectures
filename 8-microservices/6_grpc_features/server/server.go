package main

import (
	"fmt"
	"log"
	"net"
	"time"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/tap"

	"github.com/go-park-mail-ru/lectures/8-microservices/5_grpc/session"
)

// request -> {InTapHandle} | parsing | -> {Interceptor} | handling | ...

const authLogFormat = `--
	after incoming call=%v
	req=%#v
	reply=%#v
	time=%v
	md=%v
	err=%v
`

func authInterceptor(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (interface{}, error) {
	start := time.Now()

	md, _ := metadata.FromIncomingContext(ctx)

	// ... auth logic

	reply, err := handler(ctx, req)

	fmt.Printf(authLogFormat, info.FullMethod, req, reply, time.Since(start), md, err)

	return reply, err
}

func rateLimiter(ctx context.Context, info *tap.Info) (context.Context, error) {
	fmt.Printf("--\ncheck ratelimiter for %s\n", info.FullMethodName)

	// ... rate limit logic

	return ctx, nil
}

func main() {
	lis, err := net.Listen("tcp", ":8081")
	if err != nil {
		log.Fatalln("can't listen port", err)
	}

	server := grpc.NewServer(
		grpc.UnaryInterceptor(authInterceptor),
		grpc.InTapHandle(rateLimiter),
	)

	session.RegisterAuthCheckerServer(server, NewSessionManager())

	fmt.Println("starting server at :8081")
	server.Serve(lis)
}
