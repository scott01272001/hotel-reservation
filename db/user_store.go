package db

import (
	"context"
	"fmt"
	"reflect"

	"github.com/scott/hotel-reservation/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const userColl = "users"

type Dropper interface {
	Drop(context.Context) error
}

type UserStore interface {
	Dropper

	GetUserById(context.Context, string) (*types.User, error)
	GetUsers(context.Context) (*[]types.User, error)
	InsertUser(context.Context, *types.User) (*types.User, error)
	DeleteUser(context.Context, string) error
	UpdateUser(context.Context, string, *types.UpdateUserParams) (*types.User, error)
}

func ToBson[param types.UpdateUserParams](update *param) *bson.M {
	result := bson.M{}
	val := reflect.ValueOf(update)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	typ := val.Type()
	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		if !field.IsZero() {
			jsonTag := typ.Field(i).Tag.Get("json")
			if jsonTag != "" {
				result[jsonTag] = field.Interface()
			} else {
				result[typ.Field(i).Name] = field.Interface()
			}
		}
	}
	return &result
}

type MongoUserStore struct {
	client *mongo.Client
	coll   *mongo.Collection
}

func (s *MongoUserStore) Drop(ctx context.Context) error {
	fmt.Println("---- dropping users collection")
	return s.coll.Drop(ctx)
}

func (s *MongoUserStore) UpdateUser(ctx context.Context, id string, params *types.UpdateUserParams) (*types.User, error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	updateData := ToBson(params)
	update := bson.D{primitive.E{Key: "$set", Value: &updateData}}

	res, err := s.coll.UpdateByID(ctx, oid, update)
	if err != nil {
		return nil, err
	}
	if res.MatchedCount == 0 {
		return nil, fmt.Errorf("'%s' not exist", id)
	}

	var updated types.User
	if err := s.coll.FindOne(ctx, bson.M{"_id": oid}).Decode(&updated); err != nil {
		return nil, err
	}
	return &updated, nil
}

func (s *MongoUserStore) DeleteUser(ctx context.Context, id string) error {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	filter := bson.M{"_id": oid}
	res, err := s.coll.DeleteOne(ctx, filter)
	if res.DeletedCount == 0 {
		return fmt.Errorf("'%s' not exist", id)
	}
	return err
}

func (s *MongoUserStore) InsertUser(ctx context.Context, user *types.User) (*types.User, error) {
	res, err := s.coll.InsertOne(ctx, user)
	if err != nil {
		return nil, err
	}
	user.ID = res.InsertedID.(primitive.ObjectID)
	return user, nil
}

func (s *MongoUserStore) GetUserById(ctx context.Context, id string) (*types.User, error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	var user types.User
	if err := s.coll.FindOne(ctx, bson.M{"_id": oid}).Decode(&user); err != nil {
		return nil, err
	}
	return &user, nil
}

func (s *MongoUserStore) GetUsers(ctx context.Context) (*[]types.User, error) {
	cur, err := s.coll.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	var users []types.User
	if err := cur.All(ctx, &users); err != nil {
		return nil, err
	}
	return &users, nil
}

func NewMongoUserStore(client *mongo.Client) *MongoUserStore {
	return &MongoUserStore{
		client: client,
		coll:   client.Database(DBNAME).Collection(userColl),
	}
}
