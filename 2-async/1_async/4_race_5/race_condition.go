package main

import (
	"fmt"
	"sync"
	"time"
)

var balance int = 100
var mu sync.Mutex

func main() {
	// Мы запускаем две горутины.
	go withdraw("Студент 1", 100)
	go withdraw("Студент 2", 100)

	time.Sleep(1 * time.Second)
}

func withdraw(person string, amount int) {
	mu.Lock()         // Защищаем память
	defer mu.Unlock() // одновременной записи не будет

	if balance >= amount {
		fmt.Printf("%s: Ура, снимаю деньги!\n", person)
		balance -= amount
	} else {
		fmt.Printf("%s: Недостаточно средств :(\n", person)
	}
}
