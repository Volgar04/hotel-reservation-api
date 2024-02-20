package api

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
)

func ErrorHandler(c *fiber.Ctx, err error) error {
	if apiError, ok := err.(Error); ok {
		return c.Status(apiError.Code).JSON(apiError)
	}
	apiError := NewError(http.StatusInternalServerError, err.Error())
	return c.Status(apiError.Code).JSON(apiError)
}

type Error struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func (e Error) Error() string {
	return e.Message
}

func NewError(code int, message string) Error {
	return Error{
		Code:    code,
		Message: message,
	}
}

func ErrorInvalidID() Error {
	return Error{
		Code:    http.StatusBadRequest,
		Message: "invalid id given",
	}
}

func ErrorBadRequest() Error {
	return Error{
		Code:    http.StatusBadRequest,
		Message: "invalid JSON request",
	}
}

func ErrorResourceNotFound(res string) Error {
	return Error{
		Code:    http.StatusNotFound,
		Message: res + " resource not found",
	}
}

func ErrorUnauthorized() Error {
	return Error{
		Code:    http.StatusUnauthorized,
		Message: "unauthorized request",
	}
}
