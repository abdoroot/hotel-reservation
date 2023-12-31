package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/abdoroot/hotel-reservation/api"
	"github.com/abdoroot/hotel-reservation/db"
	"github.com/abdoroot/hotel-reservation/db/fixtures"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}
}

func main() {
	ctx := context.Background()
	DBNAME := os.Getenv(db.MONGODBENVDBNAME)
	DBURI := os.Getenv(db.MONGODBENVDBURI)

	log.Println("DB:", DBNAME)
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(DBURI))
	if err != nil {
		fmt.Println("db uri:", DBURI)
		panic(err)
	}
	//drop the database
	if err = client.Database(DBNAME).Drop(ctx); err != nil {
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

	//100  hotel
	for i := 1; i <= 100; i++ {
		hn := fmt.Sprintf("rand hotel %v", i)
		fixtures.AddHotel(st, hn, "Rak", 3, nil)
	}
}
