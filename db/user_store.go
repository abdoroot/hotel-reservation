package db

import (
	"context"
	"fmt"
	"os"

	"github.com/abdoroot/hotel-reservation/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	UserCollection = "users"
)

type Droper interface {
	Drop(ctx context.Context) error
}

type UserStore interface {
	Droper
	GetUserByID(context.Context, string) (*types.User, error)
	GetUserByEmail(context.Context, string) (*types.User, error)
	GetUser(context.Context, bson.M) ([]*types.User, error)
	InsertUser(context.Context, *types.User) (*types.User, error)
	DeleteUser(ctx context.Context, filter bson.M) error
	UpdateUser(ctx context.Context, filter bson.M, req bson.M) error
}

//todo implement UserStore interface using redis

type userMongoStore struct {
	client *mongo.Client
	coll   *mongo.Collection
}

func NewMongoUserStore(client *mongo.Client) *userMongoStore {
	DBName := os.Getenv(MONGODBENVDBNAME)
	return &userMongoStore{
		client: client,
		coll:   client.Database(DBName).Collection(UserCollection),
	}
}

func (m *userMongoStore) GetUserByID(ctx context.Context, id string) (*types.User, error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		fmt.Println("fail to convert user id string to objectId :", err)
		return nil, err
	}
	var user *types.User
	if err := m.coll.FindOne(ctx, bson.M{"_id": oid}).Decode(&user); err != nil {
		return nil, err
	}
	return user, nil
}

func (m *userMongoStore) GetUserByEmail(ctx context.Context, email string) (*types.User, error) {
	var user *types.User
	if err := m.coll.FindOne(ctx, bson.M{"email": email}).Decode(&user); err != nil {
		return nil, err
	}
	return user, nil
}

func (m *userMongoStore) GetUser(ctx context.Context, filter bson.M) ([]*types.User, error) {
	var users []*types.User
	cur, err := m.coll.Find(ctx, filter)
	if err != nil {
		fmt.Println("cur err", err)
		return nil, err
	}

	if err := cur.All(ctx, &users); err != nil {
		fmt.Println("error retrive data", err)
		return nil, err
	}

	return users, nil
}

func (m *userMongoStore) InsertUser(ctx context.Context, user *types.User) (*types.User, error) {
	res, err := m.coll.InsertOne(ctx, user)
	if err != nil {
		return nil, err
	}
	user.ID = res.InsertedID.(primitive.ObjectID)
	return user, nil
}

func (m *userMongoStore) UpdateUser(ctx context.Context, filter bson.M, req bson.M) error {
	update := bson.M{"$set": req}
	_, err := m.coll.UpdateOne(ctx, filter, update)
	return err
}

func (m *userMongoStore) DeleteUser(ctx context.Context, filter bson.M) error {
	res, err := m.coll.DeleteOne(ctx, filter)
	if err == nil && res.DeletedCount > 0 {
		return nil
	}
	return err
}

func (m *userMongoStore) Drop(ctx context.Context) error {
	return m.coll.Drop(ctx)
}
