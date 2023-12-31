package db

import (
	"context"
	"log"
	"os"

	"github.com/abdoroot/hotel-reservation/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	HotelCollection = "hotels"
)

type HotelStore interface {
	Droper
	InsertHotel(context.Context, *types.Hotel) (*types.Hotel, error)
	UpdatetHotel(context.Context, bson.M, bson.M) error
	GetHotels(context.Context, bson.M, Pagination) ([]*types.Hotel, error)
	GetHotelByID(context.Context, string) (*types.Hotel, error)
}

type mongoHotelStore struct {
	client *mongo.Client
	coll   *mongo.Collection
}

func NewMongoHotelStore(client *mongo.Client) *mongoHotelStore {
	DBName := os.Getenv(MONGODBENVDBNAME)
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

type Pagination struct {
	Page  int
	Limit int
}

func (s *mongoHotelStore) GetHotels(ctx context.Context, filter bson.M, pg Pagination) ([]*types.Hotel, error) {
	skip := int64((pg.Page - 1) * pg.Limit)
	limit := int64(pg.Limit)
	opts := &options.FindOptions{
		Limit: &limit,
		Skip:  &skip,
	}
	res, err := s.coll.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	hotels := []*types.Hotel{}
	if err = res.All(ctx, &hotels); err != nil {
		return nil, err
	}
	return hotels, nil
}

func (s *mongoHotelStore) GetHotelByID(ctx context.Context, id string) (*types.Hotel, error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	filter := bson.M{"_id": oid}
	hotel := &types.Hotel{}
	if err := s.coll.FindOne(ctx, filter).Decode(hotel); err != nil {
		log.Println("db err:", err)
		return nil, err
	}
	return hotel, nil
}

func (m *mongoHotelStore) Drop(ctx context.Context) error {
	return m.coll.Drop(ctx)
}
