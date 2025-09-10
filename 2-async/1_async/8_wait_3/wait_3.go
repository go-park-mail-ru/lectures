package main

import (
	"fmt"
	"strings"
	"sync"
	"time"
)

const (
	iterationsNum = 7
	goroutinesNum = 5
)

func doWork(in int) {
	for j := 0; j < iterationsNum; j++ {
		fmt.Printf(formatWork(in, j))
		time.Sleep(time.Millisecond)
	}
}

func main() {
	wg := &sync.WaitGroup{} // Инициализируем группу
	for i := 0; i < goroutinesNum; i++ {
		wg.Go(func() { // С go 1.25
			doWork(i)
		})
	}
	time.Sleep(time.Millisecond)
	wg.Wait() // Ожидаем, пока wg.Done() не приведёт счетчик к 0
}

func formatWork(in, j int) string {
	return fmt.Sprintln(strings.Repeat("  ", in), "█",
		strings.Repeat("  ", goroutinesNum-in),
		"th", in,
		"iter", j, strings.Repeat("■", j))
}
