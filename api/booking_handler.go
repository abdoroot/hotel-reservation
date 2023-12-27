package api

import (
	"fmt"

	"github.com/abdoroot/hotel-reservation/db"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
)

type bookingHandler struct {
	store *db.Store
}

func NewBookingHandler(store *db.Store) *bookingHandler {
	return &bookingHandler{
		store: store,
	}
}

func (h *bookingHandler) HandleGetbooking(c *fiber.Ctx) error {
	id := c.Params("id")
	b, err := h.store.Booking.GetBookingByID(c.Context(), id)
	if err != nil {
		return err
	}
	return c.JSON(b)
}

func (h *bookingHandler) HandleGetbookings(c *fiber.Ctx) error {
	filter := bson.M{}
	//todo refactor GetBookings to get id string and use the bussines login here
	bookings, err := h.store.Booking.GetBookings(c.Context(), filter)
	if err != nil {
		return err
	}
	return c.JSON(bookings)
}

func (h *bookingHandler) HandleGetCancelbooking(c *fiber.Ctx) error {
	id := c.Params("id")
	booking, err := h.store.Booking.GetBookingByID(c.Context(), id)
	if err != nil {
		fmt.Println("booking not found", err)
		return err
	}

	user, err := GetAuthUser(c)
	if err != nil {
		fmt.Println("user not found", err)
		return err
	}

	if user.IsAdmin && booking.UserID != user.ID {
		//admin can cancle a booking
		return fmt.Errorf("not authorized")
	}

	param := bson.M{"cancel": true}
	if err := h.store.Booking.UpdateBooking(c.Context(), id, param); err != nil {
		fmt.Println("not updated ", err)
		return fmt.Errorf("not updated")
	}
	return c.JSON(bson.M{
		"msg": "updated",
	})
}
