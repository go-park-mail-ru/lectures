package main

import (
	"fmt"
)

var idToUsername = map[int]string{
	0: "romanov",
	1: "sulaev",
	2: "dorofeev",
}

func main() {
	var id int
	for {
		_, err := fmt.Scanf("%d", &id)
		if err != nil {
			fmt.Printf("err scanf %v\n", err)
			continue
		}
		fmt.Printf("username for id %d: %s\n", id, idToUsername[id])
	}
}
