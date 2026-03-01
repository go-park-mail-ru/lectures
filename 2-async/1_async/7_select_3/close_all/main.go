package main

import (
	"fmt"
	"time"
)

func worker(id int, cancelCh chan struct{}) {
	fmt.Println("worker", id, "started")

	for {
		select {
		case <-cancelCh:
			fmt.Println("worker", id, "stopped")
			return
		default:
			time.Sleep(300 * time.Millisecond)
		}
	}
}

func main() {
	cancelCh := make(chan struct{})

	for i := 1; i <= 3; i++ {
		go worker(i, cancelCh)
	}

	time.Sleep(1 * time.Second)

	fmt.Println("sending cancel signal")
	//cancelCh <- struct{}{}
	close(cancelCh)

	time.Sleep(1 * time.Second)
}
