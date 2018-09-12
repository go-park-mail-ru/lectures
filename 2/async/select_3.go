package main

import (
	"fmt"
)

func main() {
	cancelCh := make(chan struct{})
	dataCh := make(chan int, 1)

	go func(cancelCh chan struct{}, dataCh chan int) {
		val := 0
		for {
			select {
			case dataCh <- val:
				val++
			case <-cancelCh:
				return
			}
		}
	}(cancelCh, dataCh)

	for curVal := range dataCh {
		fmt.Println("read", curVal)
		if curVal > 3 {
			fmt.Println("send cancel")
			cancelCh <- struct{}{}
			// break
		}
	}

}
