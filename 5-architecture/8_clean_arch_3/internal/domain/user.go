// Package domain содержит сущности и их правила, независимо от HTTP и БД.
package domain

import (
	"errors"
	"net/mail"
	"strings"
)

var ErrInvalidUser = errors.New("invalid user")

type User struct {
	id    int64
	name  string
	email string
}

func NewUser(name, email string) (User, error) {
	name = strings.TrimSpace(name)
	email = strings.ToLower(strings.TrimSpace(email))
	address, err := mail.ParseAddress(email)
	if name == "" || err != nil || address.Address != email {
		return User{}, ErrInvalidUser
	}
	return User{name: name, email: email}, nil
}

func RestoreUser(id int64, name, email string) (User, error) {
	if id <= 0 {
		return User{}, ErrInvalidUser
	}
	u, err := NewUser(name, email)
	if err != nil {
		return User{}, err
	}
	u.id = id
	return u, nil
}

// ChangeProfile возвращает новую версию: ошибка не меняет исходный профиль.
func (u User) ChangeProfile(name, email string) (User, error) {
	return RestoreUser(u.id, name, email)
}

func (u User) ID() int64     { return u.id }
func (u User) Name() string  { return u.name }
func (u User) Email() string { return u.email }
