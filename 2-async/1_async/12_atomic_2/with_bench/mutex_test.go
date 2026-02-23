package main

import (
	"sync"
	"sync/atomic"
	"testing"
)

var (
	totalMutex  int32
	totalAtomic int32
	mutex       = &sync.Mutex{}
)

func BenchmarkMutexParallel(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			mutex.Lock()
			totalMutex++
			mutex.Unlock()
		}
	})
}

func BenchmarkAtomicParallel(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			atomic.AddInt32(&totalAtomic, 1)
		}
	})
}

func BenchmarkLocalCounter(b *testing.B) {
	var localCounter int64
	for i := 0; i < b.N; i++ {
		localCounter++
	}
	_ = localCounter
}
