package api

import (
	"log"
	"strconv"

	"github.com/abdoroot/hotel-reservation/db"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type hotelHandler struct {
	store *db.Store
}

func NewHotelHandler(store *db.Store) *hotelHandler {
	return &hotelHandler{
		store: store,
	}
}

func (h *hotelHandler) HandleGetRooms(c *fiber.Ctx) error {
	id := c.Params("id")
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	filter := bson.M{"hotel_id": oid}
	rooms, err := h.store.Room.GetRooms(c.Context(), filter)
	if err != nil {
		return err
	}
	return c.JSON(rooms)
}

func (h *hotelHandler) HandleGethotel(c *fiber.Ctx) error {
	id := c.Params("id")
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	hotel, err := h.store.Hotel.GetHotelByID(c.Context(), oid)
	if err != nil {
		return err
	}
	return c.JSON(hotel)
}

func (h *hotelHandler) HandleGetHotels(c *fiber.Ctx) error {
	filter := bson.M{}
	//todo use c.QueryParser()
	if rooms := c.Query("rooms"); len(rooms) != 0 {
		//todo add room filter
	}
	if rating := c.Query("rating"); len(rating) != 0 {
		ratingInt, err := strconv.Atoi(rating)
		if err == nil {
			filter["rating"] = ratingInt
		}
	}
	//todo add other filter ex rating,hotel name etc ..
	hotels, err := h.store.Hotel.GetHotels(c.Context(), filter)
	if err != nil {
		log.Println(err)
		return err
	}
	return c.JSON(hotels)
}
