package main

import (
	"log"
	"sync"
)

type Book struct {
	ID    uint
	Title string
	Price uint
}

type BookStore struct {
	books []*Book
	mu    sync.RWMutex
}

func NewBookStore() *BookStore {
	return &BookStore{
		mu:    sync.RWMutex{},
		books: []*Book{},
	}
}

func (bs *BookStore) AddBook(in *Book, out *Book) error {
	log.Println("AddBook called")

	bs.mu.Lock()
	bs.books = append(bs.books, in)
	bs.mu.Unlock()
	*out = *in
	return nil
}

func (bs *BookStore) GetBooks(in int, out *[]*Book) error {
	log.Println("GetBooks called")

	bs.mu.Lock()
	defer bs.mu.Unlock()
	*out = bs.books
	return nil
}
