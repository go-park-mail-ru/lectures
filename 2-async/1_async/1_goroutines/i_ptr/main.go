package main

import (
	"fmt"
)

const goroutinesNum = 7

func main() {
	var ptrs []*int

	for i := 0; i < goroutinesNum; i++ {
		//i := i
		fmt.Println("адрес i:", &i)
		ptrs = append(ptrs, &i)
	}

	for i, pt := range ptrs {
		fmt.Println("адрес в range i:", &i)
		ptrs = append(ptrs, &i)
		fmt.Println(*pt)
	}
}
