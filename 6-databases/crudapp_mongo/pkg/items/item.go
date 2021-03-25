package items

import (
	"strconv"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Item struct {
	ID          primitive.ObjectID `json:"id" bson:"_id"`
	Title       string             `json:"title" bson:"title"`
	Description string             `json:"description" bson:"description"`
	Updated     string             `json:"updated" bson:"updated"`
}

// позволяет items handlers не импортировать sql
func (it *Item) SetUpdated(val uint32) {
	it.Updated = strconv.Itoa(int(val))
}
