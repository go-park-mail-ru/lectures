package main

import (
	"fmt"
)

var idToUsername = map[int]string{
	0: "romanov",
	1: "sulaev",
	2: "dorofeev",
}

func handling() error {
	var id int
	_, err := fmt.Scanf("%d", &id) // ask, why so many errors?
	if err != nil {
		return fmt.Errorf("failed to get username: %w", err)
	}

	username, ok := idToUsername[id]
	if !ok {
		return fmt.Errorf("no user with id %d", id)
	}

	fmt.Printf("username for id %d: %s\n", id, username)

	return nil
}

func main() {
	for {
		err := handling()
		if err != nil {
			fmt.Printf("error happened: %v\n", err)
		}
	}
}
