package main

import (
	"fmt"
	"io"

	tr "github.com/essentialkaos/translit/v3"

	"github.com/go-park-mail-ru/lectures/8-microservices/7_grpc_stream/translit"
)

type TrServer struct {
	translit.UnimplementedTransliterationServer
}

func (srv *TrServer) EnRu(inStream translit.Transliteration_EnRuServer) error {
	// Отправляем статистику клиенту:
	// var counter int
	// go func() {
	// 	for {
	// 		statTime := 5 * time.Second
	// 		inStream.Send(&translit.Word{
	// 			Word: fmt.Sprintf("stat: counter=%d for time=%s", counter, statTime),
	// 		})
	// 		counter = 0

	// 		time.Sleep(statTime)
	// 	}
	// }()

	// Обрабатываем запросы клиента:
	for {
		inWord, err := inStream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}
		out := &translit.Word{
			Word: tr.ISO9A(inWord.Word),
		}
		fmt.Println(inWord.Word, "->", out.Word)

		inStream.Send(out)

		// counter++
	}
}

func NewTr() *TrServer {
	return &TrServer{}
}
