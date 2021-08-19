package main

import (
	"fmt"
	"sync"
	"time"
)

type Task struct {
	in  int
	out *int
}

func worker(input <-chan Task) {
	for t := range input {
		// Long compute operation simulation
		time.Sleep(100 * time.Millisecond)
		*t.out = t.in * 2
	}
}

func main() {
	const (
		Workers  = 10
		Sequence = 1000
	)

	var wg sync.WaitGroup
	wg.Add(Workers)

	// Buffer equals number of workers will allow utilize workers as much as possible
	ch := make(chan Task, Workers)

	// Creates workerpool
	for i := 0; i < Workers; i++ {
		go func() {
			defer wg.Done()
			worker(ch)
		}()
	}

	res := make([]int, Sequence)

	// Produces tasks
	for i := 0; i < Sequence; i++ {
		ch <- Task{in: i, out: &res[i]}
	}

	close(ch)

	wg.Wait()

	// Prints sequence of tasks result
	for _, val := range res {
		fmt.Println(val)
	}
}
