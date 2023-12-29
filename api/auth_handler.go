package api

import (
	"fmt"
	"net/http"

	"github.com/abdoroot/hotel-reservation/db"
	"github.com/abdoroot/hotel-reservation/types"
	"github.com/gofiber/fiber/v2"
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
		return NewError(http.StatusBadRequest, "please enter valid data")
	}

	user, err := h.store.User.GetUserByEmail(c.Context(), AuthUserRequest.Email)
	if err != nil {
		fmt.Println("GetUserByEmail Err:", err) //logging
		return NewError(http.StatusBadRequest, "error email or password")
	}

	if ok := types.IsValidPassword(user.EncreptedPassword, AuthUserRequest.Password); !ok {
		return NewError(http.StatusBadRequest, "error email or password")
	}

	token, err := CreateUserJwtToken(user)
	if err != nil {
		fmt.Println("fail to create jwt token :", err) //logging
		return NewError(http.StatusBadRequest, "internal error")
	}

	resp := types.AuthResponse{
		User:  user,
		Token: token,
	}

	return c.JSON(resp)
}
