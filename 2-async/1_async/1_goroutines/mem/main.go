package main

import (
	"fmt"
	"runtime"
	"time"
)

func main() {
	var m1, m2 runtime.MemStats
	runtime.ReadMemStats(&m1)

	for i := 0; i < 10; i++ {
		go func() {
			//var b [2048]byte
			//_ = b
			//fmt.Println(b)
			//fmt.Println(i)
		}()
	}

	time.Sleep(time.Second)

	runtime.ReadMemStats(&m2)

	fmt.Printf("Allocated before: %d KB\n", m1.Alloc/1024)
	fmt.Printf("Allocated after : %d KB\n", m2.Alloc/1024)
	fmt.Printf("Per goroutine   : ~%d bytes\n", (m2.Alloc-m1.Alloc)/10)
}
