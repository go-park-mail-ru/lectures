package main

import (
	"bytes"
	"crypto/rand"
	"fmt"
	"golang.org/x/crypto/argon2"
	"reflect"
)

func hashPass(salt []byte, plainPassword string) []byte {
	hashedPass := argon2.IDKey([]byte(plainPassword), []byte(salt), 1, 64*1024, 4, 32)
	return append(salt, hashedPass...)
}

func checkPass(passHash []byte, plainPassword string) bool {
	salt := passHash[0:8]
	userPassHash := hashPass(salt, plainPassword)
	fmt.Println("salt addr: ", reflect.ValueOf(salt).UnsafePointer())
	fmt.Println("userPassHash addr: ", reflect.ValueOf(passHash).UnsafePointer())
	return bytes.Equal(userPassHash, passHash)
}

func passExample() {
	pass := "love"
	pass1 := "love1"

	// reg
	salt := make([]byte, 8)
	rand.Read(salt)
	//fmt.Printf("salt: %x\n", salt)

	fmt.Println("salt addr1: ", reflect.ValueOf(salt).UnsafePointer())
	hashedPass := hashPass(salt, pass)
	//fmt.Printf("hashedPass: %x\n", hashedPass)

	// login
	passValid := checkPass(hashedPass, pass)
	fmt.Println("hashedPass addr: ", reflect.ValueOf(hashedPass).UnsafePointer())
	fmt.Printf("passValid: %v\n", passValid)

	passValid = checkPass(hashedPass, pass1)
	fmt.Printf("passValid: %v\n", passValid)
}

func main() {
	for i := 0; i < 3; i++ {
		fmt.Println("\titeration", i)
		passExample()
	}
}
