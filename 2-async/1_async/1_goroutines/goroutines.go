package main

import (
	"fmt"
	"time"
)

const goroutinesNum = 7

func main() {
	for i := 0; i < goroutinesNum; i++ {
		go func() { // А если go < 1.22?
			fmt.Println(i)
		}()
	}

	time.Sleep(time.Second)
}
