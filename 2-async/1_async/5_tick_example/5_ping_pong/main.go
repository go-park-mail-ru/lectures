package main

import (
	"fmt"
)

func main() {
	ch := make(chan string)

	go func() {
		for pong := range ch {
			fmt.Println(pong)
			ch <- "ping"
		}
	}()

	go func() {
		for ping := range ch {
			fmt.Println(ping)
			ch <- "pong"
		}
	}()

	ch <- "ping"

	var blockNil chan int
	<-blockNil
}
