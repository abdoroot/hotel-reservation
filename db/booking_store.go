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
	GetBookingByID(context.Context, string) (*types.Booking, error)
	GetBookings(context.Context, bson.M) ([]*types.Booking, error)
	UpdateBooking(context.Context, string, bson.M) error
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

func (s *mongoBookingStore) GetBookingByID(ctx context.Context, id string) (*types.Booking, error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	booking := &types.Booking{}
	if err := s.coll.FindOne(ctx, bson.M{"_id": oid}).Decode(&booking); err != nil {
		return nil, err
	}
	return booking, nil
}

func (s *mongoBookingStore) UpdateBooking(ctx context.Context, id string, params bson.M) error {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	update := bson.M{
		"$set": params,
	}
	if _, err := s.coll.UpdateByID(ctx, oid, update); err != nil {
		return err
	}
	return nil
}

func (s *mongoBookingStore) Drop(ctx context.Context) error {
	return s.coll.Drop(ctx)
}
