package main

import (
	"fmt"
)

func main() {
	// runtime.GOMAXPROCS(1)
	cancelCh := make(chan bool)
	dataCh := make(chan int)

	go func(cancelCh chan bool, dataCh chan int) {
		val := 0
		for {
			select {
			case <-cancelCh:
				fmt.Println("cancelled")
				close(dataCh)
				return
			case dataCh <- val:
				val++
			}
		}
	}(cancelCh, dataCh)

	for curVal := range dataCh {
		fmt.Println("read", curVal)
		if curVal > 3 {
			fmt.Println("send cancel")
			cancelCh <- true // будет ждать отработки в select
			// break
		}
	}

}
