package main

import (
	"fmt"

	prsonModule "visibility/person"
)

func init() {
	fmt.Println("I am init in main")
}

func main() {
	person := 1
	fmt.Println("I am main", person, prsonModule.Person{})
}
