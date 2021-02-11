package models

import "errors"

type User struct {
	ID       uint32
	Login    string
	Password string
}

var (
	ErrNoUser  = errors.New("no user found")
	ErrBadPass = errors.New("invalid password")
)
