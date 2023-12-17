package db

import (
	"github.com/abdoroot/hotel-reservation/types"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	UserCollection = "users"
	mongodbName    = "hotel-reservation"
)

type UserStore interface {
	GetUserByID(*fiber.Ctx, primitive.ObjectID) (*types.User, error)
	GetUsers(*fiber.Ctx) ([]*types.User, error)
	InsertUser(*fiber.Ctx, *types.User) (*types.User, error)
	DeleteUser(ctx *fiber.Ctx, filter bson.M) error
	UpdateUser(ctx *fiber.Ctx, filter bson.M, req bson.M) error
}

type mongoStore struct {
	client *mongo.Client
	coll   *mongo.Collection
}

func NewMongoStore(client *mongo.Client) *mongoStore {
	mdb := client.Database(mongodbName)
	return &mongoStore{
		client: client,
		coll:   mdb.Collection(UserCollection),
	}
}

func (m *mongoStore) GetUserByID(ctx *fiber.Ctx, id primitive.ObjectID) (*types.User, error) {
	var user *types.User
	if err := m.coll.FindOne(ctx.Context(), bson.M{"_id": id}).Decode(&user); err != nil {
		return nil, err
	}
	return user, nil
}

func (m *mongoStore) GetUsers(ctx *fiber.Ctx) ([]*types.User, error) {
	var users []*types.User
	cur, err := m.coll.Find(ctx.Context(), bson.M{})
	if err != nil {
		return nil, err
	}

	if err := cur.All(ctx.Context(), &users); err != nil {
		return nil, err
	}

	return users, nil
}

func (m *mongoStore) InsertUser(ctx *fiber.Ctx, user *types.User) (*types.User, error) {
	res, err := m.coll.InsertOne(ctx.Context(), user)
	if err != nil {
		return nil, err
	}
	user.ID = res.InsertedID.(primitive.ObjectID)
	return user, nil
}

func (m *mongoStore) UpdateUser(ctx *fiber.Ctx, filter bson.M, req bson.M) error {
	update := bson.M{"$set": req}
	_, err := m.coll.UpdateOne(ctx.Context(), filter, update)
	return err
}

func (m *mongoStore) DeleteUser(ctx *fiber.Ctx, filter bson.M) error {
	res, err := m.coll.DeleteOne(ctx.Context(), filter)
	if err == nil && res.DeletedCount > 0 {
		return nil
	}
	return err
}
