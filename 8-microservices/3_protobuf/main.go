package main

import (
	"encoding/json"
	"fmt"

	"github.com/golang/protobuf/proto"
	"gopkg.in/vmihailenco/msgpack.v2"
)

// protoc --go_out=. *.proto

func main() {
	sess := &Session{
		Login:     "vasiliy",
		Useragent: "Chrome",
	}

	dataJson, _ := json.Marshal(sess)

	fmt.Printf("dataJson: %s\n", string(dataJson))
	fmt.Printf("dataJson\nlen %d\n%v\n", len(dataJson), dataJson)

	/*
		39 байт
		{"login":"dmitry","useragent":"Chrome"}
	*/

	dataPb, _ := proto.Marshal(sess)
	fmt.Printf("dataPb\nlen %d\n%v\n", len(dataPb), dataPb)

	/*
		17 байт
		[10 7 114 118 97 115 105 108 121 18 6 67 104 114 111 109 101]

			10 // номер поля + тип
			7  // длина данных
				114 118 97 115 105 108 121

			18 // номер поля + тип
			6  // длина данных
				67 104 114 111 109 101
	*/

	dataMsgPack, _ := msgpack.Marshal(sess)
	fmt.Printf("dataMsgPack\nlen %d\n%v\n", len(dataMsgPack), dataMsgPack)
}
