package string

import (
	"regexp"
	"strings"
	"testing"
)

var (
	browser = "Mozilla/5.0 (Windows NT 6.3; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.110 Safari/537.36"
	re      = regexp.MustCompile("Android")
)

// regexp.MatchString компилирует регулярку каждый раз
func BenchmarkRegExp(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, _ = regexp.MatchString("Android", browser)
	}
}

// используем прекомпилированную регулярку
func BenchmarkRegCompiled(b *testing.B) {
	for i := 0; i < b.N; i++ {
		re.MatchString(browser)
	}
}

// просто ищем вхождение подстроки
func BenchmarkStrContains(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = strings.Contains(browser, "Android")
	}
}

/*
	go test -bench . string/string_test.go
	go test -bench . -benchmem string/string_test.go
	go test -bench . -benchmem -benchtime=100000x -cpuprofile=cpu.out -memprofile=mem.out -memprofilerate=1 string/string_test.go

	go tool pprof -http=:8083 string.test cpu.out
	go tool pprof string.test cpu.out
	go tool pprof string.test mem.out

	go tool pprof -svg -inuse_space string.test mem.out > mem_is.svg
	go tool pprof -svg -inuse_objects string.test mem.out > mem_io.svg
	go tool pprof -svg string.test cpu.out > cpu.svg

	go tool pprof -png string.test cpu.out > cpu.png
*/
