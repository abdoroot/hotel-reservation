package fixtures

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/abdoroot/hotel-reservation/db"
	"github.com/abdoroot/hotel-reservation/types"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// - for test and seed perpose

func AddUser(store *db.Store, fname, lname string, isAdmin bool) *types.User {
	ctx := context.Background()
	userreq := types.CreateUserRequest{
		FirstName:         fname,
		LastName:          lname,
		Email:             fmt.Sprintf("%v@%v.com", fname, lname),
		EncreptedPassword: fmt.Sprintf("%v_%v", fname, lname),
	}

	user, err := userreq.CreateUserFromUserRequest()
	user.IsAdmin = isAdmin
	if err != nil {
		log.Fatal(err)
	}

	if _, err := store.User.InsertUser(ctx, user); err != nil {
		log.Fatal(err)
	}
	return user
}

func AddHotel(store *db.Store, hotelname, location string, rating int, rooms []primitive.ObjectID) *types.Hotel {
	ctx := context.Background()
	var rm []primitive.ObjectID
	if rooms == nil {
		rm = []primitive.ObjectID{}
	} else {
		rm = rooms
	}
	hotel := types.Hotel{
		Name:     hotelname,
		Location: location,
		Rooms:    rm,
		Rating:   rating,
	}

	insertedHotel, err := store.Hotel.InsertHotel(ctx, &hotel)
	if err != nil {
		panic(err)
	}
	return insertedHotel
}

func AddRoom(store *db.Store, hotelId primitive.ObjectID, size string, price float64, seaSide bool) *types.Room {
	ctx := context.Background()
	room := &types.Room{
		HotelID: hotelId,
		Size:    size,
		Price:   price,
		SeaSide: seaSide,
	}
	insertedRoom, err := store.Room.InsertRoom(ctx, room)
	if err != nil {
		fmt.Printf("seed room error %v", err)
	}
	return insertedRoom
}

func AddBooking(store *db.Store, userId, roomId primitive.ObjectID, fromDate, tillDate time.Time) *types.Booking {
	b := &types.Booking{
		UserID:   userId,
		RoomID:   roomId,
		FromDate: fromDate,
		TillDate: tillDate,
	}

	ib, err := store.Booking.InsertBooking(context.TODO(), b)
	if err != nil {
		fmt.Println("error instering booking", err)
	}
	return ib
}
