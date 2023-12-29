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
		return NewError(http.StatusBadRequest, "bad request")
	}
	bookingParam := &types.BookingParam{}
	if err := c.BodyParser(bookingParam); err != nil {
		return NewError(http.StatusBadRequest, "bad request")
	}

	if err = bookingParam.Validate(); err != nil {
		return NewError(http.StatusBadRequest, fmt.Sprintf("invalid data: %v", err))
	}

	//auth user
	user, ok := c.Context().Value("user").(*types.User)
	if !ok {
		return NewError(http.StatusInternalServerError, "internal error")
	}

	//check for room availability
	ok, err = h.isRoomAvailable(c.Context(), roomID, bookingParam)
	if err != nil {
		return err
	}

	if !ok {
		return ErrorReourceNotFound(roomID.Hex())
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
		return ErrorInternalErr()
	}

	return c.JSON(insertedBooking)
}

func (h *roomHandler) isRoomAvailable(ctx context.Context, roomID primitive.ObjectID, param *types.BookingParam) (bool, error) {
	//filter by chatgpt :D
	filter := bson.M{
		"room_id": roomID,
		"$or": []bson.M{
			bson.M{
				"from_date": bson.M{"$gte": param.FromDate, "$lte": param.TillDate},
			},
			bson.M{
				"till_date": bson.M{"$gte": param.FromDate, "$lte": param.TillDate},
			},
			bson.M{
				"from_date": bson.M{"$lte": param.FromDate},
				"till_date": bson.M{"$gte": param.TillDate},
			},
		},
	}
	booking, err := h.store.Booking.GetBookings(ctx, filter)
	if err != nil {
		fmt.Println(err)
		return false, ErrorReourceNotFound("booking")
	}

	ok := len(booking) == 0
	return ok, nil
}
