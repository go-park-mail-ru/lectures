package main

import (
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/pkg/errors"
)

type MyOwnError struct {
	Reason string
	Code   int
}

func (e MyOwnError) Error() string {
	return fmt.Sprintf("error with code %d happened: %s", e.Code, e.Reason)
}

var (
	ownError1 error = MyOwnError{
		Reason: "really bad things happened",
		Code:   42,
	}
)

func someJob() error {
	if rand.Intn(100) > 50 {
		return ownError1
	}

	if randomInt := rand.Intn(100); randomInt > 50 {
		return MyOwnError{
			Reason: "random was above 50",
			Code:   randomInt,
		}
	}

	return nil
}

func jobWrapper() error {
	err := someJob()
	return errors.Wrap(err, "failed to do some job")
}

func main() {
	rand.Seed(time.Now().UnixNano())

	err := someJob()
	// if err == ownError1 // WRONG
	if errors.Is(err, ownError1) {
		log.Println("found ownError1")
		return
	}

	var myOwnError MyOwnError
	if ok := errors.As(err, &myOwnError); ok {
		log.Println("found my own error sturct")
		return
	}

	if err != nil {
		log.Println("well, that was unexpected")
		return
	}

	log.Println("error is nil")
}
