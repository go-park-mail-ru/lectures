package main

import (
	"encoding/json"
	"fmt"
)

var jsonStr = `[
	{"id": 17, "username": "iivan", "phone": 0},
	{"id": "17", "address": "none", "company": "Mail.ru"},
	5,
	[]
]`

func main() {
	data := []byte(jsonStr)

	var users interface{}
	json.Unmarshal(data, &users)
	fmt.Printf("unpacked in empty interface:\n%#v\n\n", users)

	user2 := map[string]interface{}{
		"id":       42,
		"username": "rvasily",
	}
	var user2i interface{} = user2
	result, _ := json.Marshal(user2i)
	fmt.Printf("json string from map:\n %s\n", string(result))
}
