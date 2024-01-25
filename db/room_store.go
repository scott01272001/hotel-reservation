package db

import (
	"context"

	"github.com/scott/hotel-reservation/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const roomColl = "rooms"

type RoomStore interface {
	InsertRoom(context.Context, *types.Room) (*types.Room, error)
	DeleteAll(context.Context) error
}

type MongoRoomStore struct {
	client *mongo.Client
	coll   *mongo.Collection
	HotelStore
}

func NewMongoRoomStore(client *mongo.Client, hotelStore HotelStore) *MongoRoomStore {
	return &MongoRoomStore{
		client:     client,
		coll:       client.Database(DBNAME).Collection(roomColl),
		HotelStore: hotelStore,
	}
}

func (s *MongoRoomStore) DeleteAll(ctx context.Context) error {
	_, err := s.coll.DeleteMany(ctx, bson.M{})
	return err
}

func (s *MongoRoomStore) InsertRoom(ctx context.Context, room *types.Room) (*types.Room, error) {
	session, err := s.client.StartSession()
	if err != nil {
		return nil, err
	}
	defer session.EndSession(context.TODO())

	res, err := session.WithTransaction(ctx, func(ctx mongo.SessionContext) (interface{}, error) {
		res, err := s.coll.InsertOne(ctx, room)
		if err != nil {
			return nil, err
		}
		room.ID = res.InsertedID.(primitive.ObjectID).Hex()

		filter := bson.M{"_id": room.HotelID}
		update := bson.M{"$push": bson.M{"rooms": room.ID}}
		if err := s.HotelStore.Update(ctx, filter, update); err != nil {
			return nil, err
		}

		return room, nil
	})
	if err != nil {
		return nil, err
	}
	room = res.(*types.Room)
	return room, nil
}
