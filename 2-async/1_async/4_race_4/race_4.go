package main

import (
	"fmt"
	"os"
)

func main() {
	var count int

	// Если здесь false — гонка не случится, и -race промолчит.
	enableFeature := len(os.Args) > 1 && os.Args[1] == "true"

	if enableFeature {
		// Этот код создает гонку данных, но только если мы зайдем в этот блок.
		go func() {
			count++
		}()
		count++
	}

	fmt.Println("Программа завершена. Count:", count)
}
