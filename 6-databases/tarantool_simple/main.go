package main

import (
	"fmt"

	tarantool "github.com/tarantool/go-tarantool"
)

/*
	s = box.schema.space.create('users')
	s:format({{name = 'id', type = 'unsigned'},{name = 'name', type = 'string'},{name = 'age', type = 'unsigned'}})
	s:create_index('primary', {type = 'hash', parts = {'id'}})
*/

func main() {
	opts := tarantool.Opts{User: "guest"}
	conn, err := tarantool.Connect("127.0.0.1:3301", opts)

	if err != nil {
		fmt.Println("baa: Connection refused:", err)
		return
	}

	resp, err := conn.Insert("users", []interface{}{1, "Jesus", 2019})
	if err != nil {
		fmt.Println("Error", err)
		fmt.Println("Code", resp.Code)
	}

	resp, err = conn.Select("users", "primary", 0, 1, tarantool.IterEq, []interface{}{uint(1)})
	if err != nil {
		fmt.Println("Error", err)
		fmt.Println("Code", resp.Code)
		return
	}

	for _, item := range resp.Data {
		fmt.Println(item)
	}

	resp, err = conn.Eval("return test()", []interface{}{})
	if err != nil {
		fmt.Println("Error", err)
		fmt.Println("Code", resp.Code)
		return
	}

	fmt.Println(resp.Data)

}
