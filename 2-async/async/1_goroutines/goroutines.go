package main

import (
	"fmt"
	"runtime"
	"strings"
	"time"
)

const (
	iterationsNum = 6
	goroutinesNum = 6
)

func doWork(th int) {
	for j := 0; j < iterationsNum; j++ {
		fmt.Printf(formatWork(th, j))
		// time.Sleep(time.Millisecond)
		runtime.Gosched()
	}

	go func() {
		fmt.Println("kek4")
	}()
}

func main() {
	runtime.GOMAXPROCS(0)
	for i := 0; i < goroutinesNum; i++ {
		go doWork(i)

		go func() {
			fmt.Println("kek")
		}()
	}

	time.Sleep(2 * time.Second)
}

func formatWork(in, j int) string {
	return fmt.Sprintln(strings.Repeat("  ", in), "█",
		strings.Repeat("  ", goroutinesNum-in),
		"th", in,
		"iter", j, strings.Repeat("■", j))
}

func imports() {
	fmt.Println(time.Millisecond, runtime.NumCPU())
}
