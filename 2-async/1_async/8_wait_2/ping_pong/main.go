package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	ballCh := make(chan string)
	wg := &sync.WaitGroup{}

	wg.Add(1)
	// Ping
	go func() {
		defer wg.Done()
		for {
			select {
			case ball := <-ballCh:
				fmt.Println(ball)
				time.Sleep(500 * time.Millisecond)
				ballCh <- "ping"
			}
		}
	}()

	wg.Add(1)
	// Pong
	go func() {
		defer wg.Done()
		for {
			select {
			case ball := <-ballCh:
				fmt.Println(ball)
				time.Sleep(500 * time.Millisecond)
				ballCh <- "pong"
			}
		}
	}()

	// Стартовый удар
	ballCh <- "ping"

	//select {}
	wg.Wait()
}
