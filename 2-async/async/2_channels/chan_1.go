package main

import "fmt"

func main() {
	ch1 := make(chan int, 1)

	ch1 <- 42

	go func(in chan int) {
		val := <-in // ждём записи в канал
		fmt.Println("GO: get from chan", val)
		fmt.Println("GO: after read from chan")
	}(ch1)


	ch1 <- 100500

	fmt.Println("MAIN: after put to chan")
	fmt.Scanln()
}
