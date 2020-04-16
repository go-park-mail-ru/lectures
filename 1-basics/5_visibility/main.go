package main

import (
	"fmt"

	_ "github.com/go-park-mail-ru/lectures/1/5_visibility/person"
)

func init() {
	fmt.Println("I am init in main")
}

func main() {
	fmt.Println("I am main")
}
