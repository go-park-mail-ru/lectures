package main

import (
	"sync"
	"testing"
)

func BenchmarkMapWithRWMutex(b *testing.B) {
	counters := make(map[int]int)
	var mu sync.RWMutex

	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			key := i % 1000

			if key%7 == 0 {
				mu.Lock()
				counters[key]++
				mu.Unlock()
			} else {
				mu.RLock()
				_ = counters[key]
				mu.RUnlock()
			}

			i++
		}
	})
}

func BenchmarkMapWithMutex(b *testing.B) {
	counters := make(map[int]int)
	var mu sync.Mutex

	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			key := i % 1000

			if key%7 == 0 {
				mu.Lock()
				counters[key]++
				mu.Unlock()
			} else {
				mu.Lock()
				_ = counters[key]
				mu.Unlock()
			}

			i++
		}
	})
}
