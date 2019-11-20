// go build gen/* && ./codegen.exe pack/packer.go  pack/marshaller.go
package main

import "fmt"

// lets generate code for this struct
// cgen: binpack
type User struct {
	ID       int
	RealName string `cgen:"-"`
	Name     string
	Flags    int
}

type Avatar struct {
	ID  int
	Url string
}

var test = 42

func main() {
	/*
		perl -E '$b = pack("L L/a* L", 1_123_456, "kek.kekovitch", 16);
			print map { ord.", "  } split("", $b); '
	*/
	data := []byte{
		128, 36, 17, 0,

		13, 0, 0, 0,
		107, 101, 107, 46, 107, 101, 107, 111, 118, 105, 116, 99, 104,

		16, 0, 0, 0,
	}

	u := User{}
	u.Unpack(data)
	fmt.Printf("Unpacked user %#v", u)
}
