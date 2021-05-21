package main

import (
	"fmt"
	"time"
)

func sayHello() {
	fmt.Println("Hello World")
}

func main() {
	//timer := time.AfterFunc(1*time.Second, sayHello)
	//
	//fmt.Scanln()
	//timer.Stop()
	


	timer := time.NewTimer(2*time.Second)
	t := <-timer.C

	fmt.Println("Timer", t)

	t = <-time.After(1*time.Second)

	fmt.Println("Time after", t)

}
