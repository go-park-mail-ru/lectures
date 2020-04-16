package main

import (
	"fmt"
	"reflect"
)

type UserID int

type UserID2 UserID

type User struct {
	ID       UserID2
	RealName string `unpack:"-"`
	Login    string
	Flags    int
}

func PrintReflect(u interface{}) error {
	val := reflect.ValueOf(u).Elem()

	fmt.Printf("%T have %d fields:\n", u, val.NumField())
	for i := 0; i < val.NumField(); i++ {
		valueField := val.Field(i)
		typeField := val.Type().Field(i)

		fmt.Printf("\tname=%v, type=%v, value=%v, tag=`%v`\n", typeField.Name,
			typeField.Type.Kind(),
			valueField,
			typeField.Tag.Get("unpack"),
		)
	}
	return nil
}

func main() {
	u := &User{
		ID:       42,
		RealName: "unrealname",
		Flags:    32,
	}
	err := PrintReflect(u)
	if err != nil {
		panic(err)
	}
}
