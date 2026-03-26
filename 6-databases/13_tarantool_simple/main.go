package main

import (
	"context"
	"fmt"
	"math/rand"

	"github.com/tarantool/go-tarantool/v2"
)

/*
	s = box.schema.space.create('users')
	s:format({{name = 'id', type = 'unsigned'},{name = 'name', type = 'string'},{name = 'age', type = 'unsigned'}})
	s:create_index('primary', {type = 'hash', parts = {'id'}})
*/

func main() {
	ctx := context.Background()

	dialer := tarantool.NetDialer{Address: "127.0.0.1:3301", User: "guest"}

	conn, err := tarantool.Connect(ctx, dialer, tarantool.Opts{})
	if err != nil {
		fmt.Println("baa: Connection refused:", err)
		return
	}

	resp, err := conn.Do(
		tarantool.NewInsertRequest("users").Tuple([]any{rand.Int(), fmt.Sprintf("user%d", rand.Int()), 2019}),
	).Get()
	if err != nil {
		fmt.Println("Error", err)
	}

	items := make([]any, 0)
	err = conn.Do(
		tarantool.NewSelectRequest("users").Index("primary").Offset(0).Limit(100).Iterator(tarantool.IterAll).Key([]any{uint(1)}),
	).GetTyped(&items)
	if err != nil {
		fmt.Println("Error", err)
		return
	}

	for _, item := range items {
		fmt.Println(item)
	}

	resp, err = conn.Do(
		tarantool.NewEvalRequest("return test()"),
	).Get()
	if err != nil {
		fmt.Println("Error", err)
		return
	}

	fmt.Println(resp)
}
