package db

import (
	"context"
	"log"

	"github.com/abdoroot/hotel-reservation/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	RoomsCollection = "rooms"
)

type RoomStore interface {
	Droper
	InsertRoom(context.Context, *types.Room) (*types.Room, error)
	GetRooms(context.Context, bson.M) ([]*types.Room, error)
}

type mongoRoomStore struct {
	client     *mongo.Client
	coll       *mongo.Collection
	hotelStore HotelStore
}

func NewMongoRoomStore(client *mongo.Client, hotelStore HotelStore) *mongoRoomStore {
	return &mongoRoomStore{
		client:     client,
		coll:       client.Database(DBName).Collection(RoomsCollection),
		hotelStore: hotelStore,
	}
}

func (s *mongoRoomStore) InsertRoom(ctx context.Context, room *types.Room) (*types.Room, error) {
	res, err := s.coll.InsertOne(ctx, room)
	if err != nil {
		return nil, err
	}
	insertedRoomID := res.InsertedID.(primitive.ObjectID)
	room.ID = insertedRoomID
	//Update hotel by adding the inserted room id
	fiter := bson.M{"_id": room.HotelID} //Hotel Id
	update := bson.M{"$push": bson.M{
		"rooms": insertedRoomID,
	}}
	err = s.hotelStore.UpdatetHotel(ctx, fiter, update)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return room, nil
}

func (s *mongoRoomStore) GetRooms(ctx context.Context, filter bson.M) ([]*types.Room, error) {
	resp, err := s.coll.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	rooms := []*types.Room{}
	if err := resp.All(ctx, &rooms); err != nil {
		return nil, err
	}
	return rooms, nil
}

func (m *mongoRoomStore) Drop(ctx context.Context) error {
	return m.coll.Drop(ctx)
}
