package items

import (
	"database/sql"
	"fmt"
	"strconv"

	"gorm.io/gorm"
)

type ItemZero struct {
	ID          uint32
	Title       string
	Description string
	Updated     sql.NullString
}

type Item struct {
	ID          uint32 `sql:"AUTO_INCREMENT" gorm:"primary_key"`
	Title       string
	Description string
	Updated     sql.NullString `sql:"null"`
}

// позволяет items handlers не импортировать sql
func (it *Item) SetUpdated(val uint32) {
	it.Updated = sql.NullString{String: strconv.Itoa(int(val))}
}

/*
https://gorm.io/docs/models.html
*/

func (i *Item) TableName() string {
	return "items"
}

func (i *Item) BeforeSave(*gorm.DB) (err error) {
	fmt.Println("trigger on before save")
	return
}

type Item0 struct {
	gorm.Model
	ID          uint32 `sql:"AUTO_INCREMENT" gorm:"primary_key"`
	Title       string
	Description string
	Updated     string `sql:"null"`
}
