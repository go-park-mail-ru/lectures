// go test -bench . -benchmem prealloc_test.go
package prealloc

import (
	"testing"
)

const iterNum = 1000

func BenchmarkEmptyAppend(b *testing.B) {
	for i := 0; i < b.N; i++ {
		data := make([]int, 0)
		for j := 0; j < iterNum; j++ {
			data = append(data, j)
		}
	}
}

func BenchmarkPreallocAppend(b *testing.B) {
	for i := 0; i < b.N; i++ {
		data := make([]int, 0, iterNum)
		for j := 0; j < iterNum; j++ {
			data = append(data, j)
		}
	}
}

/*
go test -bench . -benchmem -cpuprofile=cpu.out -memprofile=mem.out -memprofilerate=1 prealloc/prealloc_test.go

	go tool pprof -http=:8083 prealloc.test cpu.out
	go tool pprof prealloc.test cpu.out
	go tool pprof prealloc.test mem.out

	go tool pprof -svg -inuse_space prealloc.test mem.out > mem_is.svg
	go tool pprof -svg -inuse_objects prealloc.test mem.out > mem_io.svg
	go tool pprof -svg prealloc.test cpu.out > cpu.svg

	go tool pprof -png prealloc.test cpu.out > cpu.png


*/
