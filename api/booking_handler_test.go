package api

import (
	"encoding/json"
	"fmt"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/abdoroot/hotel-reservation/db/fixtures"
	"github.com/abdoroot/hotel-reservation/types"
	"github.com/gofiber/fiber/v2"
)

func TestAdminGetBooking(t *testing.T) {
	db := setup(t)
	defer db.tearDown(t)

	user := fixtures.AddUser(db.store, "Abdelhadi", "Moahmed", false)
	admin := fixtures.AddUser(db.store, "admin", "admin", true)
	h := fixtures.AddHotel(db.store, "Dont die while you sleeping", "Rak", 2, nil)
	r1 := fixtures.AddRoom(db.store, h.ID, "small", 99.9, false)
	booking := fixtures.AddBooking(db.store, user.ID, r1.HotelID, time.Now(), time.Now().AddDate(0, 0, 2))
	_ = booking
	app := fiber.New()
	bookingHandler := NewBookingHandler(db.store)
	adminRoute := app.Group("/", JWTAuthentication(db.store.User), AdminAuth)
	adminRoute.Get("/booking", bookingHandler.HandleGetbookings)

	req := httptest.NewRequest("GET", "/booking", nil)
	at, err := CreateUserJwtToken(admin)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Add("X-Api-Key", at)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	if resp.StatusCode != 200 {
		t.Fatalf("expexted status code 200 got %v", resp.StatusCode)
	}

	bookings := []*types.Booking{}
	if err := json.NewDecoder(resp.Body).Decode(&bookings); err != nil {
		t.Fatal(err)
	}

	if len(bookings) != 1 {
		t.Fatalf("expexted 1 booking got %v", len(bookings))
	}

	if bookings[0].ID != booking.ID {
		t.Fatal("expexted booking id to be equal")
	}

	if bookings[0].UserID != booking.UserID {
		t.Fatal("expexted booking Userid to be equal")
	}

	//Test normal user can access the booking
	req = httptest.NewRequest("GET", "/booking", nil)
	at, err = CreateUserJwtToken(user)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Add("X-Api-Key", at)
	resp, err = app.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	if resp.StatusCode == 200 {
		t.Fatalf("expexted status code non 200 ")
	}
}

func TestUserGetBooking(t *testing.T) {
	db := setup(t)
	defer db.tearDown(t)

	user := fixtures.AddUser(db.store, "Abdelhadi", "Moahmed", false)
	//nonauthuser := fixtures.AddUser(db.store, "Abdelhadi", "Moahmed", false)
	h := fixtures.AddHotel(db.store, "Dont die while you sleeping", "Rak", 2, nil)
	r1 := fixtures.AddRoom(db.store, h.ID, "small", 99.9, false)
	booking := fixtures.AddBooking(db.store, user.ID, r1.HotelID, time.Now(), time.Now().AddDate(0, 0, 2))

	app := fiber.New()
	bookingHandler := NewBookingHandler(db.store)
	authroute := app.Group("/", JWTAuthentication(db.store.User))
	authroute.Get("/booking/:id", bookingHandler.HandleGetbooking)

	req := httptest.NewRequest("GET", fmt.Sprintf("/booking/%v", booking.ID.Hex()), nil)
	token, err := CreateUserJwtToken(user)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Add("X-Api-Key", token)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	if resp.StatusCode != 200 {
		t.Fatalf("expexted status code 200 got %v", resp.StatusCode)
	}

	gotBooking := &types.Booking{}
	if err := json.NewDecoder(resp.Body).Decode(gotBooking); err != nil {
		t.Fatal(err)
	}

	if gotBooking.ID != booking.ID {
		t.Fatal("expexted booking id to be equal")
	}

	if gotBooking.UserID != booking.UserID {
		t.Fatal("expexted booking Userid to be equal")
	}
}
