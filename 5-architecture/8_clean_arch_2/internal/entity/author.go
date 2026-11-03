package entity

import (
	"errors"
	"strings"
)

var ErrInvalidAuthor = errors.New("invalid author")

type Author struct {
	id   int64
	name string
}

func NewAuthor(id int64, name string) (Author, error) {
	name = strings.TrimSpace(name)
	if id <= 0 || name == "" {
		return Author{}, ErrInvalidAuthor
	}
	return Author{id: id, name: name}, nil
}

func (a Author) ID() int64    { return a.id }
func (a Author) Name() string { return a.name }
