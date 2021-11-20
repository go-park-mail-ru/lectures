package main

import (
	"fmt"
)

func foo(lol *[]string) {
	*lol = append(*lol, "exist")
}

func boo(kek []string) {
	kek = append(kek, "destroy")
	fmt.Println("before return from function", kek)
}

func zoo(shrek []string) {
	shrek[0] = "changed"
}

func main() {
	test := []string{}
	fmt.Printf("0 create empty slice %v\n\n", test)

	fmt.Println("1 pass slice by value and append")
	boo(test)
	fmt.Println("after return from function", test)

	fmt.Printf("\n2 pass slice by ref and append\n")
	foo(&test)
	fmt.Println("after return from function", test)

	fmt.Printf("\n3 pass slice by value and change element\n")
	zoo(test)
	fmt.Println("after return from function", test)
}
