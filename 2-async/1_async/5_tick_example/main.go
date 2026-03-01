package main

import (
	"fmt"
	"time"
)

func main() {
	tick1s := time.Tick(1 * time.Second)
	tick3s := time.Tick(3 * time.Second)
	tick4s := time.Tick(4 * time.Second)
	//var blockNil chan int
	//done := make(chan struct{})

	go func() {
		for t := range tick1s {
			fmt.Println("worker-1s tick at", t.Format("15:04:05"))
		}
	}()

	go func() {
		for t := range tick3s {
			fmt.Println("worker-3s tick at", t.Format("15:04:05"))
		}
	}()

	go func() {
		for t := range tick4s {
			fmt.Println("worker-4s start", t.Format("15:04:05"))
			//close(blockNil) //panic
			//close(done) // "сигнал" для завершения
			//done <- struct{}{}
			//var blockInGo chan int
			//blockInGo <- 1
			fmt.Println("worker-4s finish", t.Format("15:04:05"))
		}
	}()

	//select {}
	//var blockNil chan int
	//<-blockNil

	/*	var ch chan string
		if ch == nil {

		}*/

	//blockNil <- 1
	//<-done
	//valCh := <-done
	//fmt.Println("valCh", valCh)
	//fmt.Println(<-done)
}
