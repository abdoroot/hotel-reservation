package api

import (
	"fmt"
	"net/http"

	"github.com/abdoroot/hotel-reservation/db"
	"github.com/abdoroot/hotel-reservation/middleware"
	"github.com/abdoroot/hotel-reservation/types"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
)

type authHandler struct {
	store *db.Store
}

func NewAuthHandler(store *db.Store) *authHandler {
	return &authHandler{
		store: store,
	}
}

// a handler should only do :
// - serilaization  of the incoming request (json)
// - do some data fetching form db
// - call some business login //Call only ******
// - return the data to user

func (h *authHandler) HandleAuthUser(c *fiber.Ctx) error {
	var AuthUserRequest types.AuthUserRequest
	if err := c.BodyParser(&AuthUserRequest); err != nil {
		fmt.Println("err parsing json body", err)
		return c.Status(http.StatusBadRequest).JSON(types.ErrorResponse{
			Msg: "please enter valid data",
		})
	}

	user, err := h.store.User.GetUserByEmail(c.Context(), AuthUserRequest.Email)
	if err != nil {
		fmt.Println("GetUserByEmail Err:", err) //logging
		return c.Status(http.StatusBadRequest).JSON(types.ErrorResponse{
			Msg: "error email or password",
		})
	}

	if ok := types.IsValidPassword(user.EncreptedPassword, AuthUserRequest.Password); !ok {
		return c.Status(http.StatusBadRequest).JSON(types.ErrorResponse{
			Msg: "error email or password",
		})
	}

	token, err := middleware.CreateUserJwtToken(user)
	if err != nil {
		fmt.Println("fail to create jwt token :", err) //logging
		return c.Status(http.StatusBadRequest).JSON(bson.M{
			"msg": fmt.Errorf("internal error"),
		})
	}

	resp := types.AuthResponse{
		User:  user,
		Token: token,
	}

	return c.JSON(resp)
}
