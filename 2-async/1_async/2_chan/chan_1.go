package main

import (
	"fmt"
)

type SomeStruct struct {
	Tmp bool
}

func main() {
	ch1 := make(chan int)

	go func(in chan int) {
		fmt.Println("GO: before read from chan")
		// time.Sleep(1000 * time.Millisecond)
		val := <-in
		fmt.Println("GO: get from chan", val)
		fmt.Println("GO: after read from chan")
	}(ch1)
	fmt.Println("MAIN: before put to chan")
	// time.Sleep(1000 * time.Millisecond)

	ch1 <- 42
	ch1 <- 100500

	fmt.Println("MAIN: after put to chan")
	fmt.Scanln()
}
