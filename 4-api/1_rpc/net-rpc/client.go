package main

import (
	"log"
	"net/rpc"
)

func main() {

	client, err := rpc.DialHTTP("tcp", "localhost:8081")
	if err != nil {
		log.Fatal("dialing:", err)
	}

	res := new(Book)
	client.Call("BookStore.AddBook", &Book{Title: "The Moon is a harsh mistress"}, res)
	if err != nil {
		log.Printf("AddBook error: %s\n", err)
	}
	log.Printf("AddBook: %#v", res)

	books := &[]*Book{}

	err = client.Call("BookStore.GetBooks", 0, books)
	if err != nil {
		log.Printf("GetBooks error: %s\n", err)
	}

	log.Printf("GetBooks: %#v", *books)
}
