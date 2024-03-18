package main

import "fmt"

type Person struct {
	Id   int
	Name string
}

func SumAges(p1, p2 Person) int {
	return p1.Id + p2.Id
}

// не изменит оригинальной структуры, для который вызван
func (p Person) UpdateName(name string) {
	p.Name = name
}

// изменяет оригинальную структуру
func (p *Person) SetName(name string) {
	p.Name = name
}

type Account struct {
	Id   int
	Name string
	Person
}

func (p *Account) SetName(name string) {
	p.Name = name
}

type MySlice []int

func (sl *MySlice) Add(val int) {
	*sl = append(*sl, val)
}

func (sl *MySlice) Count() int {
	return len(*sl)
}

func main() {
	pers := Person{Id: 1, Name: "Vasily"}
	// pers := new(Person)
	pers.SetName("Vasily Romanov")
	// (&pers).SetName("Vasily Romanov")
	// fmt.Printf("updated person: %#v\n", pers)

	var acc Account = Account{
		Id:   1,
		Name: "rvasily",
		Person: Person{
			Id:   2,
			Name: "Vasily Romanov",
		},
	}

	acc.SetName("romanov.vasily")
	acc.Person.SetName("Test")

	// fmt.Printf("%#v \n", acc)

	sl := MySlice([]int{1, 2})
	sl.Add(5)
	fmt.Println(sl.Count(), sl)
}
