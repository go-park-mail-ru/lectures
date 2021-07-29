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

func startWorker(in int, waiter *sync.WaitGroup) {
	defer waiter.Done() // wait_2.go уменьшаем счетчик на 1
	localWg := &sync.WaitGroup{}
	for j := 0; j < iterationsNum; j++ {
		fmt.Printf(formatWork(in, j))
		time.Sleep(time.Millisecond) // попробуйте убрать этот sleep
	}

	localWg.Wait()
}

func main() {
	wg := &sync.WaitGroup{} // wait_2.go инициализируем группу
	for i := 0; i < goroutinesNum; i++ {
		// wg.Add надо вызывать в той горутине, которая порождает воркеров
		// иначе другая горутина может не успеть запуститься и выполнится Wait
		wg.Add(1) // wait_2.go добавляем
		go startWorker(i, wg)
	}
	wg.Wait() // 0 wait_2.go ожидаем, пока waiter.Done() не приведёт счетчик к 0

	fmt.Println(11111)

}

func formatWork(in, j int) string {
	return fmt.Sprintln(strings.Repeat("  ", in), "█",
		strings.Repeat("  ", goroutinesNum-in),
		"th", in,
		"iter", j, strings.Repeat("■", j))
}
