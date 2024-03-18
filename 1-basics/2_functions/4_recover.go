package main

import (
	"fmt"
)

func deferTest() {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("panic happend SECOND:", err)
		}

		fmt.Println("Second defer")
	}()
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("panic happend FIRST:", err)
			panic("second panic")
		}
	}()
	fmt.Println("Some userful work")
	panic("something bad happend")
	return
}

func main() {
	deferTest()

	fmt.Println("kek")

	return
}
