package items

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type RepoMongo struct {
	Coll *mongo.Collection
}

func NewMongoRepository(coll *mongo.Collection) *RepoMongo {
	return &RepoMongo{Coll: coll}
}

func (repo *RepoMongo) GetAll() ([]*Item, error) {
	items := make([]*Item, 0, 10)
	res, err := repo.Coll.Find(context.Background(), bson.M{})
	err = res.All(context.Background(), &items)
	if err != nil {
		return nil, err
	}
	return items, nil
}

func (repo *RepoMongo) GetByID(id string) (*Item, error) {
	post := &Item{}
	oid, _ := primitive.ObjectIDFromHex(id)
	err := repo.Coll.FindOne(context.Background(), bson.M{"_id": oid}).Decode(&post)
	if err != nil {
		return nil, err
	}
	return post, nil
}

func (repo *RepoMongo) Add(ctx context.Context, elem *Item) (string, error) {
	elem.ID = primitive.NewObjectID()
	_, err := repo.Coll.InsertOne(ctx, elem)
	if err != nil {
		return "", err
	}
	// log.Println(res, elem)
	return elem.ID.Hex(), nil
}

func (repo *RepoMongo) Update(elem *Item) (int64, error) {

	res, err := repo.Coll.UpdateOne(context.Background(),
		bson.M{"_id": elem.ID},
		bson.M{
			"$set": bson.M{
				"title":       elem.Title,
				"description": elem.Description,
				"updated":     "rvasily",
			},
		},
	)
	if err != nil {
		return 0, err
	}
	return res.ModifiedCount, nil
}

func (repo *RepoMongo) Delete(id string) (int64, error) {
	oid, _ := primitive.ObjectIDFromHex(id)
	res, err := repo.Coll.DeleteOne(context.Background(), bson.M{"_id": oid})
	if err != nil {
		return 0, err
	}
	return res.DeletedCount, nil
}
