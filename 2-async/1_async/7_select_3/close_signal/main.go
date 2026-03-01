package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func worker(id int, done <-chan struct{}) {
	for {
		select {
		case <-done:
			fmt.Println("worker", id, "stopped")
			return
		default:
			fmt.Println("worker", id, "working")
			time.Sleep(1 * time.Second)
		}
	}
}

func main() {
	done := make(chan struct{})

	go worker(1, done)
	go worker(2, done)

	// канал для сигналов ОС
	sigCh := make(chan os.Signal)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)

	// ждём Ctrl+C
	<-sigCh
	fmt.Println("\nCtrl+C pressed")

	// оповещаем все горутины
	// можем обновить конфиг и т.д.
	close(done)

	time.Sleep(1 * time.Second) //подождем вывод
	fmt.Println("main exit")
}
