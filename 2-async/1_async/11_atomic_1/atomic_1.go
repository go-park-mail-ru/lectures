package main

import (
	"fmt"
	"sync"
	"time"
)

var totalOperations int32 = 0
var mu = &sync.Mutex{}

func inc() {
	mu.Lock()
	defer mu.Unlock()

	// не атомарная операция
	totalOperations++
}

func main() {
	// runtime.GOMAXPROCS(1)
	for i := 0; i < 1000; i++ {
		go inc()
	}
	time.Sleep(100 * time.Millisecond)
	// ождается 1000
	mu.Lock()
	fmt.Println("total operation = ", totalOperations)
	mu.Unlock()
}
