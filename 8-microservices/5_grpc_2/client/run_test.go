package main

import (
	"context"
	"log"
	"testing"

	"github.com/go-park-mail-ru/lectures/8-microservices/5_grpc_2/session"
	"google.golang.org/grpc"
)

func TestRun(t *testing.T) {

	grcpConn, err := grpc.Dial(
		"127.0.0.1:8081",
		grpc.WithInsecure(),
	)
	if err != nil {
		log.Fatalf("cant connect to grpc")
	}
	defer grcpConn.Close()

	sessManager := session.NewAuthCheckerClient(grcpConn)

	ctx := context.Background()

	// создаем сессию
	sessId, err := sessManager.Create(ctx,
		&session.Session{
			Login:     "rvasily",
			Useragent: "chrome",
		})
	log.Println("sessId", sessId, err)

	// проеряем сессию
	sess, err := sessManager.Check(ctx,
		&session.SessionID{
			ID: sessId.ID,
		})
	log.Println("sess", sess, err)

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
	log.Println("sess", sess, err)
	t.Fail()
}
