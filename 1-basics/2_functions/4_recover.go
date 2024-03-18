package main

import (
	"fmt"
)

func tmp() {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("panic happend SECOND:", err)
		}

		fmt.Println("Second defer")
	}()

	fmt.Println("Some userful work")

	fmt.Println("Some userful after panic")
}

func main() {
	tmp()
	fmt.Println("panic heppened")

	return
}
