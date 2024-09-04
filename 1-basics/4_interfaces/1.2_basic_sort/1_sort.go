package main

import (
	"fmt"
	"sort"
)

type Student struct {
	Age  int
	Name string
}

type Students []Student

func (sts Students) Len() int           { return len(sts) }
func (sts Students) Swap(a, b int)      { sts[a], sts[b] = sts[b], sts[a] }
func (sts Students) Less(a, b int) bool { return sts[a].Age < sts[b].Age }

func main() {
	students := Students{
		{Age: 18, Name: "Harry"},
		{Age: 76, Name: "Benjamin"},
		{Age: 20, Name: "Bob"},
		{Age: 19, Name: "Alice"},
	}

	sort.Sort(students)

	for _, s := range students {
		fmt.Printf("%d %s\n", s.Age, s.Name)
	}
}
