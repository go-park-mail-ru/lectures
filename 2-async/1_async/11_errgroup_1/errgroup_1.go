package main

import (
	"fmt"
	"golang.org/x/sync/errgroup"
	"time"
)

const (
	goroutinesNum = 3
)

func scanAndPrintInt() error {
	var i int
	_, err := fmt.Scanf("%d", &i)
	if err != nil {
		return err
	}

	fmt.Println(i)
	return nil
}

func main() {
	eg := &errgroup.Group{} // Инициализируем группу
	for i := 0; i < goroutinesNum; i++ {
		eg.Go(scanAndPrintInt)
	}
	time.Sleep(time.Millisecond)
	err := eg.Wait()
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("done")
	}
}
