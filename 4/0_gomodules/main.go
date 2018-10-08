package main

import (
	"fmt"

	"modulestest/internal"

	"github.com/satori/go.uuid"
)

func main() {
	fmt.Println(uuid.NewV4())

	internal.Internal()
}
