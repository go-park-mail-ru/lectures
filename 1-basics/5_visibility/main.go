package main

import (
	"fmt"

	// "github.com/go-park-mail-ru/lectures/1-basics/5_visibility/person"
	"github.com/go-park-mail-ru/lectures/1-basics/5_visibility/person"
)

func init() {
	fmt.Println("I am init in main")
}

func main() {
	personOne := 1
	person.Person{}
	fmt.Println("I am main", personOne, person.Person{})
}
