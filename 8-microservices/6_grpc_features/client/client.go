package main

import (
	"fmt"
	"log"
	"time"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"

	"github.com/go-park-mail-ru/lectures/8-microservices/5_grpc/session"
)

// {Interceptor} | serialization | -> request

const timingLogFormat = `--
	call=%v
	req=%#v
	reply=%#v
	time=%v
	err=%v
`

func timingInterceptor(
	ctx context.Context,
	method string,
	req interface{},
	reply interface{},
	cc *grpc.ClientConn,
	invoker grpc.UnaryInvoker,
	opts ...grpc.CallOption,
) error {
	start := time.Now()
	err := invoker(ctx, method, req, reply, cc, opts...)

	fmt.Printf(timingLogFormat, method, req, reply, time.Since(start), err)

	return err
}

type tokenAuth struct {
	Token string
}

func (t *tokenAuth) GetRequestMetadata(context.Context, ...string) (map[string]string, error) {
	return map[string]string{
		"access-token": t.Token,
	}, nil
}

func (c *tokenAuth) RequireTransportSecurity() bool {
	return false
}

func main() {
	grpcClient, err := grpc.NewClient(
		"127.0.0.1:8081",
		grpc.WithUnaryInterceptor(timingInterceptor),
		grpc.WithPerRPCCredentials(&tokenAuth{"100500"}),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatalf("cant connect to grpc")
	}
	defer grpcClient.Close()

	sessManager := session.NewAuthCheckerClient(grpcClient)

	ctx := context.Background()
	md := metadata.Pairs(
		"api-req-id", "123",
		"subsystem", "cli",
	)
	ctx = metadata.NewOutgoingContext(ctx, md)

	var header, trailer metadata.MD

	// создаем сессию
	sessionID, err := sessManager.Create(ctx,
		&session.Session{
			Login:     "rvasily",
			Useragent: "chrome",
		},
		grpc.Header(&header),
		grpc.Trailer(&trailer),
	)
	fmt.Println("sessId", sessionID, err)
	fmt.Println("header", header)
	fmt.Println("trailer", trailer)

	// проверяем сессию
	sess, err := sessManager.Check(ctx,
		&session.SessionID{
			ID: sessionID.ID,
		})
	fmt.Println("sess", sess, err)

	// удаляем сессию
	_, err = sessManager.Delete(ctx,
		&session.SessionID{
			ID: sessionID.ID,
		})

	// проверяем еще раз
	sess, err = sessManager.Check(ctx,
		&session.SessionID{
			ID: sessionID.ID,
		})
	fmt.Println("sess", sess, err)
}
