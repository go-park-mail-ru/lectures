package main

import (
	"fmt"
	"runtime"
	"strings"
	"sync"
)

const (
	iterationsNum = 6
	goroutinesNum = 6
)

func doWork(wg *sync.WaitGroup, th int) {
	defer wg.Done()
	for j := 0; j < iterationsNum; j++ {
		fmt.Printf(formatWork(th, j))
		runtime.Gosched()
	}
}

func main() {
	wg := &sync.WaitGroup{}
	runtime.GOMAXPROCS(1)
	wg.Add(goroutinesNum)
	for i := 0; i < goroutinesNum; i++ {
		go doWork(wg, i)
	}

	wg.Wait()

}

func formatWork(in, j int) string {
	return fmt.Sprintln(strings.Repeat("  ", in), "█",
		strings.Repeat("  ", goroutinesNum-in),
		"th", in,
		"iter", j, strings.Repeat("■", j))
}
