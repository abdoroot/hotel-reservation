package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/abdoroot/hotel-reservation/db"
	"github.com/abdoroot/hotel-reservation/middleware"
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
	us     db.UserStore
	bs     db.BookingStore
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
	us = db.NewMongoUserStore(client)
	bs = db.NewMongoBookingStore(client)
}

func seedHotel(hotelname, location string, rating int) *types.Hotel {
	hotel := types.Hotel{
		Name:     hotelname,
		Location: location,
		Rooms:    []primitive.ObjectID{},
		Rating:   rating,
	}

	insertedHotel, err := hs.InsertHotel(ctx, &hotel)
	if err != nil {
		panic(err)
	}
	return insertedHotel
}

func seedRooms(hotelId primitive.ObjectID, size string, price float64, seaSide bool) *types.Room {
	room := &types.Room{
		HotelID: hotelId,
		Size:    size,
		Price:   price,
		SeaSide: seaSide,
	}
	insertedRoom, err := rs.InsertRoom(ctx, room)
	if err != nil {
		fmt.Printf("seed room error %v", err)
	}
	return insertedRoom
}

func seedUser(fname, lname, email string, isAdmin bool) *types.User {
	userreq := types.CreateUserRequest{
		FirstName:         fname,
		LastName:          lname,
		Email:             email,
		EncreptedPassword: "abdoroot123",
	}

	user, err := userreq.CreateUserFromUserRequest()
	user.IsAdmin = isAdmin
	if err != nil {
		log.Fatal(err)
	}

	if _, err := us.InsertUser(context.TODO(), user); err != nil {
		log.Fatal(err)
	}

	return user
}

func seedBooking(userId, roomId primitive.ObjectID, fromDate, tillDate time.Time) *types.Booking {
	b := &types.Booking{
		UserID:   userId,
		RoomID:   roomId,
		FromDate: fromDate,
		TillDate: tillDate,
	}

	ib, err := bs.InsertBooking(context.TODO(), b)
	if err != nil {
		fmt.Println("error instering booking", err)
	}
	return ib
}

func main() {
	u := seedUser("Abdelhadi", "Moahmed", "abd.200930@gmail.com", false)
	u2 := seedUser("admin", "admin", "admin@admin.com", true)
	h := seedHotel("Dont die while you sleeping", "Rak", 2)
	r1 := seedRooms(h.ID, "small", 99.9, false)
	seedRooms(h.ID, "normal", 120.5, false)
	seedRooms(h.ID, "kingsize", 200.99, true)
	seedBooking(u.ID, r1.HotelID, time.Now(), time.Now().AddDate(0, 0, 2))
	//get api tokens
	t, _ := middleware.CreateUserJwtToken(u)
	t2, _ := middleware.CreateUserJwtToken(u2)
	fmt.Printf("user :%v Token -> %v\n\n", u.Email, t)
	fmt.Printf("user :%v Token -> %v", u2.Email, t2)
}
