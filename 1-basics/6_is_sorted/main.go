package main

import (
	"bufio"
	"log"
	"os"

	"github.com/go-park-mail-ru/lectures/6_is_sorted/sorted"
)

func main() {
	var inputStrings []string

	in := bufio.NewScanner(os.Stdin)
	for in.Scan() {
		inputStrings = append(inputStrings, in.Text())
	}

	if err := in.Err(); err != nil {
		log.Fatalf("input scanning failed: %s", err)
	}

	if err := sorted.Check(inputStrings); err != nil {
		log.Fatalf("sorted check failed: %s", err)
	}

	log.Println("strings are sorted")
}
