package main

import "fmt"

func main() {
	// инициализация при создании
	var user map[string]string = map[string]string{
		"name":     "Anton",
		"lastName": "Sulaev",
	}

	// сразу с нужной ёмкостью
	profile := make(map[string]string, 10)

	// количество элементов
	mapLength := len(user)

	fmt.Printf("%d %+v\n", mapLength, profile)

	/*

			1. Ключ всегда comparable
			= !=

			2.
			sha1, crc32, crc64

			Fh = func (toHash Comparable, d int)  int

			Vasya -> Fh(Vasya) -> 10
			Kolya ->  Fh(Kolya) -> 10

			-----------&----- 15 адресов
					   V
					   K
			O(1)

		Может быть не конкурентно
	*/

	// если ключа нет - вернёт значение по умолчанию для типа
	mName := user["middleName"]
	fmt.Println("mName:", mName)

	// проверка на существование ключа
	mName, mNameExist := user["middleName"]
	fmt.Println("mName:", mName, "mNameExist:", mNameExist)

	// пустая переменная - только проверяем что ключ есть
	_, mNameExist2 := user["middleName"]
	fmt.Println("mNameExist2", mNameExist2)

	// удаление ключа
	delete(user, "lastName")
	fmt.Printf("%#v\n", user)
}
