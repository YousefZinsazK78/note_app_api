package api

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
)

func ErrorHandler(c *fiber.Ctx, err error) error {
	if apiErr, ok := err.(Error); !ok {
		return c.Status(apiErr.Code).JSON(apiErr.Err)
	}
	apiNewError := NewError(fiber.StatusInternalServerError, err.Error())
	return c.Status(apiNewError.Code).JSON(apiNewError.Err)
}

type Error struct {
	Code int    `json:"code"`
	Err  string `json:"err"`
}

func (e Error) Error() string {
	return e.Err
}

func NewError(code int, err string) Error {
	return Error{
		code,
		err,
	}
}

func ErrNotFound() error {
	return Error{
		Code: http.StatusNotFound,
		Err:  "resource not found!",
	}
}
