package items

import (
	"database/sql"
	"strconv"
)

type Item struct {
	ID          uint32
	Title       string
	Description string
	Updated     sql.NullString
}

// позволяет items handlers не импортировать sql
func (it *Item) SetUpdated(val uint32) {
	it.Updated = sql.NullString{String: strconv.Itoa(int(val))}
}
