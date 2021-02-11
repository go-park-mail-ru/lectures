package main

import (
	"encoding/json"
	"fmt"

	"github.com/golang/protobuf/proto"
)

func main() {
	sess := &Session{
		Login:     "dmitry",
		Useragent: "Chrome",
	}

	dataJson, _ := json.Marshal(sess)

	fmt.Printf("dataJson\nlen %d\n%v\n", len(dataJson), dataJson)

	/*
		39 байт
		{"login":"dmitry","useragent":"Chrome"}
	*/

	dataPb, _ := proto.Marshal(sess)
	fmt.Printf("dataPb\nlen %d\n%v\n", len(dataPb), string(dataPb))

	/*
		17 байт
		[10 6 100 109 105 116 114 121 18 6 67 104 114 111 109 101]

			10 // номер поля + тип
			6  // длина данных
				100 109 105 116 114 121
			18 // номер поля + тип
			6  // длина данных
				67 104 114 111 109 101
	*/

}
