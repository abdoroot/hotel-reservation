package main

import (
	"context"
	"fmt"
	"log"

	"github.com/abdoroot/hotel-reservation/db"
	"github.com/abdoroot/hotel-reservation/types"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	ctx    = context.Background()
	client *mongo.Client
	hs     db.HotelStore
	rs     db.RoomStore
)

func init() {
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(db.DBURI))
	if err != nil {
		panic(err)
	}
	//drop the database
	if err = client.Database(db.DBName).Drop(ctx); err != nil {
		log.Fatal(err)
	}

	hs = db.NewMongoHotelStore(client)
	rs = db.NewMongoRoomStore(client, hs)
}

func hotelSeed(hotelname, location string, rating int) {
	hotel := types.Hotel{
		Name:     hotelname,
		Location: location,
		Rooms:    []primitive.ObjectID{},
		Rating:   rating,
	}

	rooms := []types.Room{
		{
			Size:    "small",
			Price:   99.9,
			SeaSide: false,
		},
		{
			Size:    "normal",
			Price:   120.5,
			SeaSide: false,
		},
		{
			Size:    "kingsize",
			Price:   200.99,
			SeaSide: true,
		},
	}

	insertedHotel, err := hs.InsertHotel(ctx, &hotel)
	if err != nil {
		panic(err)
	}

	fmt.Println(insertedHotel)
	insertedRoomsIds := []primitive.ObjectID{}
	for _, val := range rooms {
		room := &val
		room.HotelID = insertedHotel.ID
		insertedRoom, err := rs.InsertRoom(ctx, room)
		if err != nil {
			fmt.Println(err)
			continue
		}
		insertedRoomsIds = append(insertedRoomsIds, insertedRoom.ID)
	}

	fmt.Println(insertedRoomsIds)
}

func main() {
	hotelSeed("Paramount Hotel", "Dubai", 4)
	hotelSeed("Dont die while you sleeping", "Rak", 2)
}
