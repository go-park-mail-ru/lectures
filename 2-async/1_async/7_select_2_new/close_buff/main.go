package main

import (
	"fmt"
)

func main() {
	ch := make(chan int, 3)

	ch <- 1
	ch <- 2
	ch <- 3

	close(ch)
	fmt.Println("channel closed")

	/*	for v := range ch {
		fmt.Println("read:", v)
	}*/
	for {
		select {
		case v, ok := <-ch:
			if !ok {
				fmt.Println("select: channel empty and closed — exit")
				return
			}
			fmt.Println("select: read value", v)
		}
	}
}
