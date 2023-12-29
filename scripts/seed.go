package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/abdoroot/hotel-reservation/api"
	"github.com/abdoroot/hotel-reservation/db"
	"github.com/abdoroot/hotel-reservation/db/fixtures"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	ctx := context.Background()
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(db.DBURI))
	if err != nil {
		panic(err)
	}
	//drop the database
	if err = client.Database(db.DBName).Drop(ctx); err != nil {
		log.Fatal(err)
	}

	hs := db.NewMongoHotelStore(client)
	rs := db.NewMongoRoomStore(client, hs)
	us := db.NewMongoUserStore(client)
	bs := db.NewMongoBookingStore(client)

	st := &db.Store{
		Hotel:   hs,
		User:    us,
		Room:    rs,
		Booking: bs,
	}

	nu := fixtures.AddUser(st, "Abdelhadi", "Moahmed", false)
	au := fixtures.AddUser(st, "admin", "admin", true)
	h := fixtures.AddHotel(st, "Dont die while you sleeping", "Rak", 2, nil)
	r1 := fixtures.AddRoom(st, h.ID, "small", 99.9, false)
	fixtures.AddRoom(st, h.ID, "normal", 120.5, false)
	fixtures.AddRoom(st, h.ID, "kingsize", 200.99, true)
	fixtures.AddBooking(st, nu.ID, r1.HotelID, time.Now(), time.Now().AddDate(0, 0, 2))
	//create api tokens
	t, _ := api.CreateUserJwtToken(nu)
	t2, _ := api.CreateUserJwtToken(au)
	fmt.Printf("user :%v Token -> %v\n\n", nu.Email, t)
	fmt.Printf("user :%v Token -> %v", au.Email, t2)
}
