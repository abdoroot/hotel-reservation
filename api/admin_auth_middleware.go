package api

import (
	"github.com/abdoroot/hotel-reservation/types"
	"github.com/gofiber/fiber/v2"
)

func AdminAuth(c *fiber.Ctx) error {
	user, ok := c.Context().Value("user").(*types.User)
	if !ok {
		return ErrorUnauthorized()
	}

	if !user.IsAdmin {
		return ErrorUnauthorized()
	}
	return c.Next()
}
