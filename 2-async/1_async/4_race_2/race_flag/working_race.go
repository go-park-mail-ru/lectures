package main

import (
	"fmt"
	"time"
)

// add race flag
func main() {
	var counter int
	for i := 0; i < 4; i++ {
		go func() {
			for j := 0; j < 1000; j++ {
				counter++
			}
		}()
	}
	time.Sleep(2 * time.Second)
	fmt.Println("counter =", counter)
}
