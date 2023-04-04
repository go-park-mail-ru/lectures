package main

import (
	"context"
	"io"
	"log"
	"sync"
	"testing"

	"github.com/go-park-mail-ru/lectures/8-microservices/6_grpc_stream/translit"
	"google.golang.org/grpc"
)

func Test(t *testing.T) {
	grcpConn, err := grpc.Dial(
		"127.0.0.1:8081",
		grpc.WithInsecure(),
	)
	if err != nil {
		log.Fatalf("cant connect to grpc")
	}
	defer grcpConn.Close()

	tr := translit.NewTransliterationClient(grcpConn)

	ctx := context.Background()
	stream, _ := tr.EnRu(ctx)

	wg := &sync.WaitGroup{}
	wg.Add(2)

	go func(wg *sync.WaitGroup) {
		defer wg.Done()
		for {
			outWord, err := stream.Recv()
			if err == io.EOF {
				log.Println("\tstream closed")
				return
			} else if err != nil {
				log.Println("\terror happed", err)
				return
			}
			log.Println(" <-", outWord.Word)
		}
	}(wg)

	go func(wg *sync.WaitGroup) {
		defer wg.Done()
		words := []string{"privet", "kak", "dela"}
		for _, w := range words {
			log.Println("-> ", w)
			stream.Send(&translit.Word{
				Word: w,
			})
			//time.Sleep(2 * time.Second)
		}
		stream.CloseSend()
		log.Println("\tsend done")
	}(wg)

	wg.Wait()

	t.Fail()
}
