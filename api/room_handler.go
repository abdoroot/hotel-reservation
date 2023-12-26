package api

import (
	"context"
	"fmt"
	"net/http"

	"github.com/abdoroot/hotel-reservation/db"
	"github.com/abdoroot/hotel-reservation/types"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type roomHandler struct {
	store *db.Store
}

func NewRoomHandler(store *db.Store) *roomHandler {
	return &roomHandler{
		store: store,
	}
}

func (h *roomHandler) HandleRoomBooking(c *fiber.Ctx) error {
	roomID, err := primitive.ObjectIDFromHex(c.Params("id"))
	if err != nil {
		return err
	}
	bookingParam := &types.BookingParam{}
	if err := c.BodyParser(bookingParam); err != nil {
		return err
	}

	if err = bookingParam.Validate(); err != nil {
		return c.JSON(types.ErrorResponse{
			Type: "error",
			Msg:  err.Error(),
		})
	}

	//auth user
	user, ok := c.Context().Value("user").(*types.User)
	if !ok {
		return fmt.Errorf("internal error")
	}

	//check for room availability
	ok, err = h.isRoomAvailable(c.Context(), roomID, bookingParam)
	if err != nil {
		return err
	}

	if !ok {
		return c.Status(http.StatusBadRequest).JSON(types.ErrorResponse{
			Type: "error",
			Msg:  "room not avilable",
		})
	}

	//process to booking
	booking := &types.Booking{
		UserID:     user.ID,
		RoomID:     roomID,
		FromDate:   bookingParam.FromDate,
		TillDate:   bookingParam.TillDate,
		NumPersons: bookingParam.NumPersons,
	}

	insertedBooking, err := h.store.Booking.InsertBooking(c.Context(), booking)
	if err != nil {
		return err
	}

	return c.JSON(insertedBooking)
}

func (h *roomHandler) isRoomAvailable(ctx context.Context, roomID primitive.ObjectID, param *types.BookingParam) (bool, error) {
	//filter by chatgpt :D
	filter := bson.M{
		"room_id": roomID,
		"$or": []bson.M{
			// Booking starts within the given date range
			bson.M{
				"from_date": bson.M{"$gte": param.FromDate, "$lte": param.TillDate},
			},
			// Booking ends within the given date range
			bson.M{
				"till_date": bson.M{"$gte": param.FromDate, "$lte": param.TillDate},
			},
			// Booking spans the entire given date range
			bson.M{
				"from_date": bson.M{"$lte": param.FromDate},
				"till_date": bson.M{"$gte": param.TillDate},
			},
		},
	}

	booking, err := h.store.Booking.GetBookings(ctx, filter)
	if err != nil {
		fmt.Println(err)
		return false, fmt.Errorf("room not avilable")
	}

	ok := len(booking) == 0
	return ok, nil
}
