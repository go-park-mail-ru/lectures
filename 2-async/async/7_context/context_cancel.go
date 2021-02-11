package main

import (
	"context"
	"fmt"
	"math/rand"
	"time"
)

func студент(ctx context.Context, workerNum int, out chan<- int) {
	waitTime := time.Duration(rand.Intn(100)+10) * time.Millisecond
	fmt.Println(workerNum, "студент думает", waitTime)
	select {
	case <-ctx.Done():
		return
	case <-time.After(waitTime):
		fmt.Println("студент", workerNum, "придумал")
		out <- workerNum
	}
}

func main() {
	ctx, finish := context.WithCancel(context.Background())
	result := make(chan int, 1)

	for i := 0; i <= 10; i++ {
		go студент(ctx, i, result)
	}

	foundBy := <-result
	fmt.Println("вопрос был задан студентом", foundBy)
	finish()

	time.Sleep(time.Second)
}
