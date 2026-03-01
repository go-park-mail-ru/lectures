package main

import (
	"errors"
	"fmt"
)

type MyError struct {
	Code int
}

func (e MyError) Error() string {
	return fmt.Sprintf("my error, code=%d", e.Code)
}

func bad() error {
	var err *MyError = nil
	fmt.Println("in bad", err == nil)
	return nil // явно указываем nil
}

func main() {
	err := bad()
	var me MyError
	fmt.Println("errors as: ", errors.As(err, &me))
	fmt.Println("in main", err == nil)
	fmt.Printf("err: %#v\n", err)
}
