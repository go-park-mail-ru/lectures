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

func doWork(in int, wg *sync.WaitGroup) {
	defer wg.Done() // wait_2.go уменьшаем счетчик на 1
	for j := 0; j < iterationsNum; j++ {
		fmt.Printf(formatWork(in, j))
		time.Sleep(time.Millisecond)
	}
}

func main() {
	wg := &sync.WaitGroup{} // wait_2.go инициализируем группу
	for i := 0; i < goroutinesNum; i++ {
		// wg.Add надо вызывать в той горутине, которая порождает воркеров
		// иначе другая горутина может не успеть запуститься и выполнится Wait
		wg.Add(1) // wait_2.go добавляем
		go doWork(i, wg)
	}
	time.Sleep(time.Millisecond)
	wg.Wait() // wait_2.go ожидаем, пока wg.Done() не приведёт счетчик к 0
	// fmt.Scanln(a ...any)

	// fmt.Println(11111)

}

func formatWork(in, j int) string {
	return fmt.Sprintln(strings.Repeat("  ", in), "█",
		strings.Repeat("  ", goroutinesNum-in),
		"th", in,
		"iter", j, strings.Repeat("■", j))
}
