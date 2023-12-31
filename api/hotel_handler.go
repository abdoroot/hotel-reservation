package api

import (
	"log"
	"net/http"

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
		return NewError(http.StatusBadRequest, "bad request")
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
	hotel, err := h.store.Hotel.GetHotelByID(c.Context(), id)
	if err != nil {
		return NewError(http.StatusBadRequest, "bad request")
	}
	return c.JSON(hotel)
}

type HotelParam struct {
	db.Pagination
	Rating int
}

type DataResponse struct {
	Result int `json:"result"`
	Data   any `json:"data"`
	Page   int `json:"page"`
}

func (h *hotelHandler) HandleGetHotels(c *fiber.Ctx) error {
	filter := bson.M{}
	params := HotelParam{}
	if err := c.QueryParser(&params); err != nil {
		ErrorBadRequest()
	}

	if params.Rating != 0 {
		filter["rating"] = params.Rating
	}

	hotels, err := h.store.Hotel.GetHotels(c.Context(), filter, params.Pagination)
	if err != nil {
		log.Println(err)
		return NewError(http.StatusInternalServerError, "internal error")
	}
	resp := DataResponse{
		Result: len(hotels),
		Data:   hotels,
		Page:   params.Pagination.Page,
	}
	return c.JSON(resp)
}
