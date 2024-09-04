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
			panic(err)
		}

		username, ok := idToUsername[id]
		if !ok {
			panic(fmt.Sprintf("no user with id %d", id))
		}

		fmt.Printf("username for id %d: %s\n", id, username)
	}
}
