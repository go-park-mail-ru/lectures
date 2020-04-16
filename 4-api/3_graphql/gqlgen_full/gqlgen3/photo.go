package gqlgen3

import (
	// "log"
	"strconv"
)

type Photo struct {
	ID     uint `json:"id"`
	UserID uint `json:"-"`
	// User     *User  `json:"user"`
	URL     string `json:"url"`
	Comment string `json:"comment"`
	Rating  int    `json:"rating"`
	Liked   bool   `json:"liked"`
}

func (ph *Photo) Id() string {
	// log.Println("call Photo.Id method", ph.ID)
	return strconv.Itoa(int(ph.ID))
}
