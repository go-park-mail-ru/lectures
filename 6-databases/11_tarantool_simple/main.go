package main

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"github.com/tarantool/go-tarantool/v2"
)

/*
	s = box.schema.space.create('users')
	s:format({{name = 'id', type = 'unsigned'},{name = 'name', type = 'string'},{name = 'age', type = 'unsigned'}})
	s:create_index('primary', {type = 'hash', parts = {'id'}})
*/

func main() {
	ctx := context.Background()

	rand.Seed(time.Now().UnixNano())

	dialer := tarantool.NetDialer{Address: "127.0.0.1:3301", User: "guest"}

	conn, err := tarantool.Connect(ctx, dialer, tarantool.Opts{})

	if err != nil {
		fmt.Println("baa: Connection refused:", err)
		return
	}

	resp, err := conn.Insert("users", []interface{}{rand.Int(), fmt.Sprintf("user%d", rand.Int()), 2019})
	if err != nil {
		fmt.Println("Error", err)
	}

	resp, err = conn.Select("users", "primary", 0, 100, tarantool.IterAll, []interface{}{uint(1)})
	if err != nil {
		fmt.Println("Error", err)
		return
	}

	for _, item := range resp {
		fmt.Println(item)
	}

	resp, err = conn.Eval("return test()", []interface{}{})
	if err != nil {
		fmt.Println("Error", err)
		return
	}

	fmt.Println(resp)

}
