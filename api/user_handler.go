package api

import (
	"errors"

	"github.com/abdoroot/hotel-reservation/db"
	"github.com/abdoroot/hotel-reservation/types"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type userHandler struct {
	store *db.Store
}

func NewUserHandler(store *db.Store) *userHandler {
	return &userHandler{
		store: store,
	}
}

func (h *userHandler) HandleGetUser(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	user, err := h.store.User.GetUserByID(ctx.Context(), oid)
	if err != nil {
		return err
	}
	return ctx.JSON(user)
}

func (h *userHandler) HandlePostUser(c *fiber.Ctx) error {
	var userRequest types.CreateUserRequest
	if err := c.BodyParser(&userRequest); err != nil {
		return c.JSON("err")
	}
	if errs := userRequest.Validate(); len(errs) > 0 {
		return errors.Join(errs...)
	}

	user, err := userRequest.CreateUserFromUserRequest()
	if err != nil {
		return err
	}

	res, err := h.store.User.InsertUser(c.Context(), user)
	if err != nil {
		return err
	}

	return c.JSON(res)
}

func (h *userHandler) HandlePutUser(c *fiber.Ctx) error {
	updateRequest := &types.UpdateRequest{}
	if err := c.BodyParser(updateRequest); err != nil {
		return err
	}
	if errs := updateRequest.Validate(); len(errs) != 0 {
		return errors.Join(errs...)
	}

	id := c.Params("id")
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	filter := bson.M{
		"_id": oid,
	}
	return h.store.User.UpdateUser(c.Context(), filter, updateRequest.ToBSON())
}

func (h *userHandler) HandleGetUsers(c *fiber.Ctx) error {
	filter := bson.M{}
	users, err := h.store.User.GetUser(c.Context(), filter)
	if err != nil {
		return err
	}
	return c.JSON(users)
}

func (h *userHandler) HandleDeleteUser(c *fiber.Ctx) error {
	id := c.Params("id")
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	filter := bson.M{"_id": oid}
	err = h.store.User.DeleteUser(c.Context(), filter)
	if err == nil {
		return c.JSON("deleted")
	}
	return err
}
