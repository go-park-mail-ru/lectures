package main

import (
	"fmt"
)

/*
	go run -gcflags -m main.go
	go run -gcflags '-m -m' main.go
*/

type User struct {
	ID int
	Login string
}

func (u *User) GetID() int {
	return u.ID
}

func newUser(login string) *User {
	return &User{123, login}
}

func setToZero(in *int) {
	// for i := 0; i<3; i++ {
	// 	*in = 1
	// }
	*in = 0
}

func main() {

	u := newUser("test")
	u.ID = 1

	data := make([]string, 20)
	data = append(data, "test")

	i := 1
	setToZero(&i)

	// _ = fmt.Sprint(data)
	// _ = fmt.Sprint(u)
	fmt.Println("test")
}
