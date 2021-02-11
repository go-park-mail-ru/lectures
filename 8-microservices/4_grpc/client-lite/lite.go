package main

import (
	"context"
	"fmt"
	"log"

	"google.golang.org/grpc"

	"github.com/go-park-mail-ru/lectures/8-microservices/4_grpc/session"
)

func main() {

	grpcConn, err := grpc.Dial(
		"127.0.0.1:8081",
		grpc.WithInsecure(),
	)
	if err != nil {
		log.Fatalf("cant connect to grpc")
	}
	defer grpcConn.Close()

	sessManager := session.NewAuthCheckerClient(grpcConn)

	ctx := context.Background()

	// создаем сессию
	sessId, err := sessManager.Create(ctx,
		&session.Session{
			Login:     "rvasily",
			Useragent: "chrome",
		})
	fmt.Println("sessId", sessId, err)

	// проеряем сессию
	sess, err := sessManager.Check(ctx,
		&session.SessionID{
			ID: sessId.ID,
		})
	fmt.Println("sess", sess, err)

	// удаляем сессию
	_, err = sessManager.Delete(ctx,
		&session.SessionID{
			ID: sessId.ID,
		})

	// проверяем еще раз
	sess, err = sessManager.Check(ctx,
		&session.SessionID{
			ID: sessId.ID,
		})
	fmt.Println("sess", sess, err)
}
