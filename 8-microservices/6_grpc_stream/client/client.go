package main

import (
	"context"
	"fmt"
	"github.com/go-park-mail-ru/lectures/8-microservices/6_grpc_stream/translit"
	"google.golang.org/grpc"
	"io"
	"log"
	"sync"
)

func main() {

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
	stream, err := tr.EnRu(ctx)

	wg := &sync.WaitGroup{}
	wg.Add(2)

	go func(wg *sync.WaitGroup) {
		defer wg.Done()
		words := []string{"privet", "kak", "dela", "go", "vyhodit'", "s", "udalyonki"}
		for _, w := range words {
			fmt.Println("-> ", w)
			stream.Send(&translit.Word{
				Word: w,
			})
			//time.Sleep(1 * time.Second)
		}
		stream.CloseSend()
		fmt.Println("\tsend done")
	}(wg)

	go func(wg *sync.WaitGroup) {
		defer wg.Done()
		for {
			outWord, err := stream.Recv()
			if err == io.EOF {
				fmt.Println("\tstream closed")
				return
			} else if err != nil {
				fmt.Println("\terror happed", err)
				return
			}
			fmt.Println(" <-", outWord.Word)
		}
	}(wg)

	wg.Wait()

}
