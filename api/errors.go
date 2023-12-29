package api

import (
	"fmt"
	"net/http"

	"github.com/gofiber/fiber/v2"
)

type Error struct {
	Status int    `json:"status"`
	Msg    string `json:"message"`
}

func (e Error) Error() string {
	return e.Msg
}

func NewError(status int, msg string) Error {
	return Error{
		Status: status,
		Msg:    msg,
	}
}

func ErrorBadRequest() Error {
	return Error{
		Status: http.StatusBadRequest,
		Msg:    "Bad request",
	}
}

func ErrorInternalErr() Error {
	return Error{
		Status: http.StatusInternalServerError,
		Msg:    "inernal error",
	}
}

func ErrorReourceNotFound(resource string) Error {
	return Error{
		Status: http.StatusBadRequest,
		Msg:    fmt.Sprintf("%v resource not found", resource),
	}
}

func ErrorUnauthorized() Error {
	return Error{
		Status: http.StatusUnauthorized,
		Msg:    "Un Authorized",
	}
}

func ErrorHandler(ctx *fiber.Ctx, err error) error {
	apiErr, ok := err.(Error)
	if ok {
		return ctx.Status(apiErr.Status).JSON(err)
	}
	newApiErr := NewError(http.StatusInternalServerError, err.Error())
	return ctx.Status(newApiErr.Status).JSON(newApiErr)
}
