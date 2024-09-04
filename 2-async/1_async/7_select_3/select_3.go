package main

import (
	"fmt"
)

func main() {
	cancelCh := make(chan bool)
	dataCh := make(chan int)

	go func(cancelCh chan bool, dataCh chan int) {
		val := 0
		for {
			select {
			case <-cancelCh:
				// закрваем канал там где пишем
				// а тут что если много генераторов ??
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
			cancelCh <- true
			// что если этот канал должен остановить множество генераторов ?
			// close(cancelCh)
		}
	}
}
