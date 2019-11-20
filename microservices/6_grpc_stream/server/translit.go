package main

import (
	"fmt"
	"io"
	"time"

	"github.com/go-park-mail-ru/lectures/microservices/6_grpc_stream/translit"

	tr "github.com/gen1us2k/go-translit"
)

type TrServer struct {
}

func (srv *TrServer) EnRu(inStream translit.Transliteration_EnRuServer) error {
	for {
		time.Sleep(1 * time.Second)
		inWord, err := inStream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}
		out := &translit.Word{
			Word: tr.Translit(inWord.Word),
		}
		fmt.Println(inWord.Word, "->", out.Word)
		inStream.Send(out)
	}
	return nil
}

func NewTr() *TrServer {
	return &TrServer{}
}
