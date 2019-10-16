package model

type User struct {
	Username string
	Password string `json:"-"`
}
