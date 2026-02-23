package main

import (
	"fmt"
	"time"
)

func main() {
	ticker := time.NewTicker(time.Second)
	i := 0
	for tickTime := range ticker.C {
		i++
		fmt.Println("step", i, "time", tickTime)
		if i >= 5 {
			// Надо останавливать, иначе потечет (но с Go 1.23 не надо 🙂)
			ticker.Stop()
			//os.Exit(3)
			break
		}
	}
	fmt.Println("total", i)

	return

	// Не может быть остановлен и собран сборщиком мусора
	// Используйте, если должен работать вечно
	c := time.Tick(time.Second)
	i = 0
	for tickTime := range c {
		i++
		fmt.Println("step", i, "time", tickTime)
		if i >= 5 {
			break
		}
	}

}

//while true; do echo "Starting program..."; go run tick.go; echo "Program exited with code $?"; sleep 1; done;
