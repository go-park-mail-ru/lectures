package main

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"sync"

	"google.golang.org/grpc/credentials/insecure"

	"google.golang.org/grpc"

	"github.com/go-park-mail-ru/lectures/8-microservices/7_grpc_stream/translit"
)

func main() {
	grpcClient, err := grpc.NewClient(
		"127.0.0.1:8081",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatalf("failed to connect to grpc: %v", err)
	}
	defer grpcClient.Close()

	tr := translit.NewTransliterationClient(grpcClient)

	ctx := context.Background()
	stream, err := tr.EnRu(ctx)
	if err != nil {
		log.Fatalf("failed to get stream: %v", err)
	}

	wg := &sync.WaitGroup{}
	wg.Add(1)

	// Вычитываем данные из стрима:
	go func(wg *sync.WaitGroup) {
		defer wg.Done()
		for {
			outWord, err := stream.Recv()
			if err == io.EOF {
				fmt.Println("\tstream closed")
				return
			} else if err != nil {
				fmt.Println("\terror happened", err)
				return
			}

			fmt.Println(" <-", outWord.Word)
		}
	}(wg)

	// Читаем со стандартного ввода и отправляем в стрим:
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		err = stream.Send(&translit.Word{
			Word: scanner.Text(),
		})
		if err != nil {
			fmt.Printf("\tfailed to send data to stream: %v\n", err)
		}
	}

	stream.CloseSend()

	wg.Wait()
}
