package main

import (
	"fmt"

	"github.com/go-park-mail-ru/lectures/8-microservices/9_flatbuffers/session"
	flatbuffers "github.com/google/flatbuffers/go"
)

// flatc --go *.fbs

// Можно и для других языков:
// flatc --cpp *.fbs
// flatc --java *.fbs
// flatc --python *.fbs

func main() {
	// Сериализация:
	builder := flatbuffers.NewBuilder(1024) // Создаём сериализатор

	// Сериализуем вложенные структуры - []int64:
	session.SessionStartIdsVector(builder, 3) // Начинаем сериализацию Ids
	builder.PrependInt64(42)
	builder.PrependInt64(100500)
	builder.PrependInt64(999)
	ids := builder.EndVector(3) // Заканчиваем сериализацию Ids

	// Сериализуем вложенные структуры - string:
	ua := builder.CreateString("mozilla")
	name := builder.CreateString("bob")

	session.SessionStart(builder) // Начинаем сериализацию Session
	session.SessionAddLogin(builder, name)
	session.SessionAddUserAgent(builder, ua)
	session.SessionAddIds(builder, ids)
	s := session.SessionEnd(builder) // Заканчиваем сериализацию Session

	builder.Finish(s) // Заканчиваем сериализацию

	result := builder.FinishedBytes()

	/*
		80 байт
		[16 0 0 0 0 0 10 0 16 0 12 0 8 0 4 0 10 0 0 0 32 0 0 0 16 0 0 0 4 0 0 0 3 0 0 0 98 111 98 0 7 0 0 0 109 111 122 105 108 108 97 0 3 0 0 0 231 3 0 0 0 0 0 0 148 136 1 0 0 0 0 0 42 0 0 0 0 0 0 0]

		98 111 98 - bob
		109 111 122 105 108 108 97 - mozilla
	*/

	fmt.Println("flatbuffers:")
	fmt.Printf("\tlen: %d\n", len(result))
	fmt.Printf("\traw data: %v\n", result)
	fmt.Println()

	// Десериализация:
	sDeserialized := session.GetRootAsSession(result, 0)

	fmt.Println("deserialized session:")
	fmt.Printf("\tlogin = %s\n", sDeserialized.Login())
	fmt.Printf("\tids[1] = %d\n", sDeserialized.Ids(1))
	fmt.Printf("\tids[0] = %d\n", sDeserialized.Ids(0))
	fmt.Printf("\tids[2] = %d\n", sDeserialized.Ids(2))
	fmt.Printf("\tuser agent = %s\n", sDeserialized.UserAgent())
}
