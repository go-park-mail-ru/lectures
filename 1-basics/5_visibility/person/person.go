package person

import "fmt"

var (
	Public  = 1
	private = 1
)

func init() {
	fmt.Println("I am init 1")
}

func init() {
	fmt.Println("I am init 2")
}

type Person struct {
	ID     int
	Name   string
	secret string
}

func (p Person) UpdateSecret(secret string) {
	fmt.Println(private)
	p.secret = secret
}
