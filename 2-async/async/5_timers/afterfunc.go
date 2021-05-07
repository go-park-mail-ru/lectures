package main

import (
	"fmt"
	"sync"
	"time"
)

const idle = 3

type Resource struct {
	mu  sync.Mutex
	count int
	timer *time.Timer
}

func (r *Resource) Update() int {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.timer != nil && !r.timer.Stop() {
		r.timer.Reset(idle*time.Second)
		return r.count
	}
	r.timer = time.AfterFunc(idle*time.Second, func() {
		r.mu.Lock()
		defer r.mu.Unlock()

		r.count += 1
		r.timer.Stop()
	})
	return r.count
}

func main() {
	r := Resource{}
	// Opens new connection
	r.Update()
	fmt.Printf("conn count -> %d\n", r.count)

	// Takes old connection. Do not increment count
	r.Update()
	fmt.Printf("conn count -> %d\n", r.count)

	// Takes old connection. Do not increment count
	r.Update()
	fmt.Printf("conn count -> %d\n", r.count)

	// Uncomment this line to imitate idle timeout and use new connection
	// time.Sleep((idle + 1) * time.Second)

	r.Update()
	fmt.Printf("conn count -> %d\n", r.count)
}
