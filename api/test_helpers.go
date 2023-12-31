package api

import (
	"context"
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/abdoroot/hotel-reservation/db"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type testdb struct {
	store *db.Store
}

func init() {
	if err := godotenv.Load("../.env"); err != nil {
		log.Fatal("Error loading .env file")
	}
}

func setup(t *testing.T) *testdb {
	DBURI := os.Getenv(db.MONGODBENVDBURI)
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(DBURI))
	if err != nil {
		panic(err)
	}

	hs := db.NewMongoHotelStore(client)
	rs := db.NewMongoRoomStore(client, hs)
	us := db.NewMongoUserStore(client)
	bs := db.NewMongoBookingStore(client)
	store := &db.Store{
		Hotel:   hs,
		User:    us,
		Room:    rs,
		Booking: bs,
	}
	return &testdb{
		store: store,
	}
}

func (tdb *testdb) tearDown(t *testing.T) {
	ctx := context.TODO()
	fmt.Println("--- dropping collection")
	tdb.store.User.Drop(ctx)
	tdb.store.Hotel.Drop(ctx)
	tdb.store.Room.Drop(ctx)
	tdb.store.Booking.Drop(ctx)
}
