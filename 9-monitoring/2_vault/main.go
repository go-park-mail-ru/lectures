package main

import (
	"fmt"
	"log"
	"os"

	"github.com/hashicorp/vault/api"
)

func main() {
	client, err := api.NewClient(&api.Config{
		Address: "http://127.0.0.1:8200",
	})

	if err != nil {
		log.Fatal(err)
	}

	token := os.Getenv("token")
	client.SetToken(token)
	secretValues, err := client.Logical().Read("secret/data/postgres")
	if err != nil {
		log.Println(err)
	}
	fmt.Printf("secret %s -> %+v", "postgres", secretValues.Data)
}
