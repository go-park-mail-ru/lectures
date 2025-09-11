package main

import (
	"context"
	"fmt"
	"golang.org/x/sync/errgroup"
	"time"
)

const (
	goroutinesNum  = 3
	badGorutineNum = 2
)

func printGorutineNum(ctx context.Context, num int) error {
	waitTime := time.Duration(10*(num+1)) * time.Millisecond
	fmt.Println(num, "gorutine wil work after", waitTime)

	select {
	case <-ctx.Done():
		fmt.Printf("gorutine %d cancelled\n", num)
		return nil
	case <-time.After(waitTime):
		if num == badGorutineNum {
			fmt.Println("error found in gorutine", num)
			return fmt.Errorf("bad gorutine number %d", num)
		}
		fmt.Println("goroutine number", num)
	}

	return nil
}

func main() {
	eg, ctx := errgroup.WithContext(context.Background()) // Инициализируем группу с контекстом
	for i := 0; i < goroutinesNum; i++ {
		eg.Go(func() error {
			return printGorutineNum(ctx, i+1)
		})
	}
	time.Sleep(time.Millisecond)
	err := eg.Wait()
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("done")
	}
}
