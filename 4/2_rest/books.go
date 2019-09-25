package main

import (
	"log"
	"sync"
)

type Book struct {
	ID    uint   `json:"id"`
	Title string `json:"title"`
	Price uint   `json:"price"`
}

type BookStore struct {
	books  []*Book
	mu     sync.RWMutex
	nextID uint
}

func NewBookStore() *BookStore {
	return &BookStore{
		mu:    sync.RWMutex{},
		books: []*Book{},
	}
}

func (bs *BookStore) AddBook(in *Book) (uint, error) {
	log.Println("AddBook called")

	bs.mu.Lock()
	bs.nextID++
	in.ID = bs.nextID
	log.Println("nextID", bs.nextID)
	bs.books = append(bs.books, in)
	bs.mu.Unlock()

	return in.ID, nil
}

func (bs *BookStore) GetBooks() ([]*Book, error) {
	log.Println("GetBooks called")

	bs.mu.RLock()
	defer bs.mu.RUnlock()

	return bs.books, nil
}
