package main

import (
	"fmt"
	"runtime"
)

func main() {
	runtime.GOMAXPROCS(1)

	in := make(chan int, 1)

	go func(out chan<- int) {
		for i := 0; i <= 10; i++ {
			fmt.Println("before", i)
			out <- i

			fmt.Println("after", i)
		}
		// out <- 12
		close(out)
		fmt.Println("generator finish")
	}(in)

	for i := range in {
		///
		fmt.Println("\tget", i)
	}

	// fmt.Scanln()
}
