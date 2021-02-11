package transaction

import (
	"reflect"
	"testing"
)

func TestList(t *testing.T) {
	transaction := Transaction{
		ID: "test",
		Amount: 10.0,
	}

	list := NewList()

	list.Add(transaction)

	if _, exists := list.transactions["test"]; !exists {
		t.Error("transaction doesn't exist")
	}

	if got := list.transactions["test"]; got != transaction {
		t.Error("transactions are not equal")
	}
}

func TestGet(t *testing.T) {
	transactionFixture := Transaction{ID: "test", Amount: 5}
	expected := []Transaction{transactionFixture}

	list := NewList()
	list.transactions["test"] = transactionFixture

	if !reflect.DeepEqual(list.Get(), expected) {
		t.Error("slices are not equal")
	}
}