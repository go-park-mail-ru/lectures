package main

import (
	"fmt"
	"runtime"
	"strings"
	"sync"
	"time"
)

const (
	iterationsNum = 6
	goroutinesNum = 5
	quotaLimit    = 2
)

func startWorker(in int, wg *sync.WaitGroup, quotaCh chan struct{}) {
	defer wg.Done()

	fmt.Printf("start %d\n", in)
	time.Sleep(time.Millisecond)

	quotaCh <- struct{}{}
	defer func() { <-quotaCh }()
	for j := 0; j < iterationsNum; j++ {
		fmt.Println(formatWork(in, j))

		// if j%2 == 0 {
		// 	<-quotaCh             // ratelim.go, возвращаем слот
		// 	quotaCh <- struct{}{} // ratelim.go, берём слот
		// }

		runtime.Gosched() // даём поработать другим горутинам
	}
}

func main() {
	wg := &sync.WaitGroup{}
	quotaCh := make(chan struct{}, quotaLimit)
	for i := 0; i < goroutinesNum; i++ {
		wg.Add(1)
		go startWorker(i, wg, quotaCh)
	}
	time.Sleep(time.Millisecond)
	wg.Wait()
}

func formatWork(in, j int) string {
	return fmt.Sprintln(strings.Repeat("  ", in), "█",
		strings.Repeat("  ", goroutinesNum-in),
		"th", in,
		"iter", j, strings.Repeat("■", j))
}
