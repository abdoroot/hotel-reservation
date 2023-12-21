package db

import (
	"context"

	"github.com/abdoroot/hotel-reservation/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	HotelCollection = "hotels"
)

type HotelStore interface {
	Droper
	InsertHotel(context.Context, *types.Hotel) (*types.Hotel, error)
	UpdatetHotel(context.Context, bson.M, bson.M) error
	GetHotel(context.Context, bson.M) ([]*types.Hotel, error)
}

type mongoHotelStore struct {
	client *mongo.Client
	coll   *mongo.Collection
}

func NewMongoHotelStore(client *mongo.Client) *mongoHotelStore {
	return &mongoHotelStore{
		client: client,
		coll:   client.Database(DBName).Collection(HotelCollection),
	}
}

func (s *mongoHotelStore) InsertHotel(ctx context.Context, hotel *types.Hotel) (*types.Hotel, error) {
	res, err := s.coll.InsertOne(ctx, hotel)
	if err != nil {
		return nil, err
	}
	hotel.ID = res.InsertedID.(primitive.ObjectID)
	return hotel, nil
}

func (s *mongoHotelStore) UpdatetHotel(ctx context.Context, filter bson.M, update bson.M) error {
	_, err := s.coll.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}
	return nil
}

func (s *mongoHotelStore) GetHotel(ctx context.Context, filter bson.M) ([]*types.Hotel, error) {
	hotels := []*types.Hotel{}
	res, err := s.coll.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	err = res.All(ctx, &hotels)
	if err != nil {
		return nil, err
	}
	return hotels, nil
}

func (m *mongoHotelStore) Drop(ctx context.Context) error {
	return m.coll.Drop(ctx)
}
