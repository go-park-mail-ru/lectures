package main

import "fmt"
import "time"

func MyFunc(in chan int) {
	fmt.Println("GO: before read from chan")
	val := <-in
	fmt.Println("GO: get from chan", val)
	fmt.Println("GO: after read from chan")
}

func main() {
	ch1 := make(chan int, 1)

	go MyFunc(ch1)

	time.Sleep(10 * time.Millisecond)

	ch1 <- 42
	ch1 <- 100500

	fmt.Println("MAIN: after put to chan")
	fmt.Scanln()
}
