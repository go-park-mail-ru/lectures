package main

import (
	"fmt"
	"sample/transaction"
	"github.com/satori/go.uuid"
	"os"
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