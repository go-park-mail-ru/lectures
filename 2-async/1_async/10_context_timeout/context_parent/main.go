package main

import (
	"context"
	"fmt"
	"time"
)

func worker(ctx context.Context, name string) {
	fmt.Println(name, "стартовал")

	select {
	case <-ctx.Done():
		fmt.Println(name, "остановлен:", ctx.Err())
	case <-time.After(3 * time.Second):
		fmt.Println(name, "успел закончить работу", ctx.Err())
	}
}

func main() {
	parentCtx, cancelParent := context.WithCancel(context.Background())
	//	parentCtx, cancelParent := context.WithTimeout(context.Background(), 500*time.Millisecond)

	childCtx1, _ := context.WithCancel(parentCtx)
	//childCtx2, _ := context.WithTimeout(parentCtx, 3*time.Second)

	go worker(childCtx1, "child-1")
	go worker(childCtx1, "child-2")

	time.Sleep(4 * time.Second)

	fmt.Println("main: отменяем родительский контекст")
	cancelParent()

	time.Sleep(1 * time.Second)
}
