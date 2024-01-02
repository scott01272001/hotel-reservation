package db

import (
	"github.com/scott/hotel-reservation/types"
	"go.mongodb.org/mongo-driver/mongo"
)

const DBNAME = "hotel-reservation"
const userColl = "users"

type UserStore interface {
	GetUserById(id string) (*types.User, error)
}

type MongoUserStore struct {
	client *mongo.Client
	coll   *mongo.Collection
}

func (s *MongoUserStore) GetUserById(id string) (*types.User, error) {
	return nil, nil
}

func NewMongoUserStore(client *mongo.Client) *MongoUserStore {
	return &MongoUserStore{
		client: client,
		coll:   client.Database(DBNAME).Collection(userColl),
	}
}
