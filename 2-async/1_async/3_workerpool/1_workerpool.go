package main

import (
	"fmt"
	"runtime"
	"strings"
	"time"
)

const goroutinesNum = 3

func startWorker(workerNum int, in <-chan string) {
	for input := range in {
		fmt.Println(formatWork(workerNum, input))
		time.Sleep(10 * time.Millisecond)
	}
	printFinishWork(workerNum)
}

func formatWork(in int, input string) string {
	return fmt.Sprintln(strings.Repeat("  ", in), "█",
		strings.Repeat("  ", goroutinesNum-in),
		"th", in,
		"recieved", input)
}

func printFinishWork(in int) {
	fmt.Println(strings.Repeat("==", in), "█",
		strings.Repeat("==", goroutinesNum-in),
		"===", in,
		"finished")
}

func main() {
	runtime.GOMAXPROCS(0)
	workerInput := make(chan string)
	for i := 0; i < goroutinesNum; i++ {
		go startWorker(i, workerInput)
	}

	months := []string{"Январь", "Февраль", "Март",
		"Апрель", "Май", "Июнь",
		"Июль", "Август", "Сентябрь",
		"Октябрь", "Ноябрь", "Декабрь",
		"123",
	}

	for _, monthName := range months {
		workerInput <- monthName
	}
	close(workerInput) // попробуйте закомментировать

	time.Sleep(time.Millisecond)
}
