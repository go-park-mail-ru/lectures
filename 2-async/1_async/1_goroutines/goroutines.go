package main

import (
	"fmt"
)

const (
	iterationsNum = 7
	goroutinesNum = 7
)

func main() {
	// memory i[ 7 ]
	for i := 0; i < goroutinesNum; i++ {
		go func(i int) {
			fmt.Println(i)
		}(i)
	}
	fmt.Scanln()
}
