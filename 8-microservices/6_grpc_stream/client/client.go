package main

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"sync"

	"google.golang.org/grpc"

	"github.com/go-park-mail-ru/lectures/8-microservices/6_grpc_stream/translit"
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
	if err != nil {
		println(err)
	}

	wg := &sync.WaitGroup{}
	wg.Add(1)

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

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		errSend := stream.Send(&translit.Word{
			Word: scanner.Text(),
		})
		if err != nil {
			println(errSend)
		}
	}
	stream.CloseSend()

	wg.Wait()

}
