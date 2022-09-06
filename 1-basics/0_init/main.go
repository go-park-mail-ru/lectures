package main

import (
	"fmt"
	"os"

	"github.com/go-park-mail-ru/lectures/1-basics/0_init/transaction"
	uuid "github.com/satori/go.uuid"
)

func main() {
	fmt.Println("hi")

	list := transaction.NewList()

	t := transaction.Transaction{
		ID: uuid.NewV4().String(),
		Account: transaction.Account{
			Title: "Sberbank",
		},
		Amount: 10000,
	}

	list.Add(t)

	fmt.Println(list.Get())

	fmt.Println(os.Args)
}
