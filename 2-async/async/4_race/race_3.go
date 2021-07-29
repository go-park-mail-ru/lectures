package main

import (
	"fmt"
	"sync"
)

func main() {
	var counters = &sync.Map{}
	for i := 0; i < 5; i++ {
		go func(counters *sync.Map, th int) {
			for j := 0; j < 5; j++ {
				counters.Store(th*10+j, 5)
			}
		}(counters, i)
	}

	fmt.Scanln()
	fmt.Println("counters result", counters)
}
