package main

import (
	"encoding/json"
	"fmt"

	"google.golang.org/protobuf/proto"

	msgpack "github.com/vmihailenco/msgpack/v5"
)

// protoc --go_out=. *.proto

// Можно и для других языков:
// protoc --cpp_out=. *.proto
// protoc --java_out=. *.proto
// protoc --js_out=. *.proto

func main() {
	session := &Session{
		Login:     "vasiliy",
		Useragent: "Chrome",
		Ids:       []int64{42, 100500},
	}

	dataJson, _ := json.Marshal(session)

	fmt.Println("json:")
	fmt.Printf("\tstring data: %s\n", string(dataJson))
	fmt.Printf("\tlen: %d\n\traw data: %v\n", len(dataJson), dataJson)

	/*
		58 байт
		{"login":"vasiliy","useragent":"Chrome","ids":[42,100500]}
	*/

	dataPB, _ := proto.Marshal(session)
	fmt.Println("protobuf:")
	fmt.Printf("\tlen: %d\n\traw data: %v\n", len(dataPB), dataPB)

	/*
		23 байта
		[10 7 118 97 115 105 108 105 121 18 6 67 104 114 111 109 101 26 4 42 148 145 6]

			10 // номер поля + тип
			7  // длина данных
				118 97 115 105 108 105 121 - vasiliy

			18 // номер поля + тип
			6  // длина данных
				67 104 114 111 109 101 - Chrome

			26 // номер поля + тип
				4 42 148 145 6 - [42,100500]
	*/

	dataMsgPack, _ := msgpack.Marshal(session)
	fmt.Println("msgpack:")
	fmt.Printf("\tstring data: %s\n", string(dataMsgPack))
	fmt.Printf("\tlen: %d\n\traw data: %v\n", len(dataMsgPack), dataMsgPack)

	/*
		55 байт
		[131 165 76 111 103 105 110 167 118 97 115 105 108 105 121 169 85 115 101 114 97 103 101 110 116 166 67 104 114 111 109 101 163 73 100 115 146 211 0 0 0 0 0 0 0 42 211 0 0 0 0 0 1 136 148]
	*/
}
