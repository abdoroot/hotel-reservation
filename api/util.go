package api

import (
	"fmt"

	"github.com/abdoroot/hotel-reservation/types"
	"github.com/gofiber/fiber/v2"
)

func GetAuthUser(c *fiber.Ctx) (*types.User, error) {
	user, ok := c.Context().Value("user").(*types.User)
	if !ok {
		return nil, fmt.Errorf("not authorized")
	}
	return user, nil
}
