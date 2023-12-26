package db

import (
	"context"

	"github.com/abdoroot/hotel-reservation/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	bookingCollection = "booking"
)

type BookingStore interface {
	Droper
	InsertBooking(context.Context, *types.Booking) (*types.Booking, error)
	GetBookings(context.Context,bson.M) ([]*types.Booking, error)
}

type mongoBookingStore struct {
	client *mongo.Client
	coll   *mongo.Collection
}

func NewMongoBookingStore(client *mongo.Client) *mongoBookingStore {
	return &mongoBookingStore{
		client: client,
		coll:   client.Database(DBName).Collection(bookingCollection),
	}
}

func (s *mongoBookingStore) InsertBooking(ctx context.Context, booking *types.Booking) (*types.Booking, error) {
	res, err := s.coll.InsertOne(ctx, booking)
	if err != nil {
		return nil, err
	}
	booking.ID = res.InsertedID.(primitive.ObjectID)
	return booking, nil
}

func (s *mongoBookingStore) GetBookings(ctx context.Context, filter bson.M) ([]*types.Booking, error) {
	cur, err := s.coll.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	bookings := []*types.Booking{}
	if err = cur.All(ctx, &bookings); err != nil {
		return nil, err
	}
	return bookings, nil
}

func (s *mongoBookingStore) Drop(ctx context.Context) error {
	return s.coll.Drop(ctx)
}
