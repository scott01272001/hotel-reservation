package db

import (
	"context"

	"github.com/scott/hotel-reservation/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const hotelColl = "hotels"

type HotelStore interface {
	InsertHotel(context.Context, *types.Hotel) (*types.Hotel, error)
	Update(context.Context, bson.M, bson.M) error
	DeleteAll(context.Context) error
	GetHotels(context.Context, bson.M) (*[]types.Hotel, error)
	GetRoomsByHotelId(context.Context, string) (*[]types.Room, error)
}

type MongoHotelStore struct {
	client *mongo.Client
	coll   *mongo.Collection
}

func NewMongoHotelStore(client *mongo.Client) *MongoHotelStore {
	return &MongoHotelStore{
		client: client,
		coll:   client.Database(DBNAME).Collection(hotelColl),
	}
}

func (s *MongoHotelStore) GetRoomsByHotelId(ctx context.Context, id string) (*[]types.Room, error) {
	hotelID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	pipeline := mongo.Pipeline{
		bson.D{{"$match", bson.D{{"_id", hotelID}}}},
		bson.D{{"$lookup", bson.D{
			{"from", "rooms"},
			{"let", bson.D{{"roomIds", "$rooms"}}},
			{"pipeline", bson.A{bson.D{
				{"$match", bson.D{{"$expr", bson.D{{"$in", bson.A{"$_id", bson.D{{"$map", bson.D{
					{"input", "$$roomIds"},
					{"as", "roomId"},
					{"in", bson.D{{"$toObjectId", "$$roomId"}}},
				}}}}}}}}},
			}}},
			{"as", "roomDetails"},
		}}},
		bson.D{{"$unwind", "$roomDetails"}},
		bson.D{{"$replaceRoot", bson.D{{"newRoot", "$roomDetails"}}}},
	}
	cur, err := s.coll.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	var rooms []types.Room
	if err = cur.All(context.TODO(), &rooms); err != nil {
		return nil, err
	}
	return &rooms, nil
}

func (s *MongoHotelStore) GetHotels(ctx context.Context, filter bson.M) (*[]types.Hotel, error) {
	cur, err := s.coll.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	var hotels []types.Hotel
	if err := cur.All(ctx, &hotels); err != nil {
		return nil, err
	}
	return &hotels, nil
}

func (s *MongoHotelStore) DeleteAll(ctx context.Context) error {
	_, err := s.coll.DeleteMany(ctx, bson.M{})
	return err
}

func (s *MongoHotelStore) Update(ctx context.Context, filter bson.M, update bson.M) error {
	_, err := s.coll.UpdateOne(ctx, filter, update)
	return err
}

func (s *MongoHotelStore) InsertHotel(ctx context.Context, hotel *types.Hotel) (*types.Hotel, error) {
	res, err := s.coll.InsertOne(ctx, hotel)
	if err != nil {
		return nil, err
	}
	hotel.ID = res.InsertedID.(primitive.ObjectID).Hex()
	return hotel, nil
}
