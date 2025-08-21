package handler

import (
	"github.com/budsx/expenses-management/model"
	"github.com/gofiber/fiber/v2"
)

func SuccessResponse(c *fiber.Ctx, message string, data interface{}) error {
	return c.Status(fiber.StatusOK).JSON(model.Response{
		Message: message,
		Data:    data,
	})
}

func ErrorResponse(c *fiber.Ctx, statusCode int, errorType string, message string) error {
	return c.Status(statusCode).JSON(model.ErrorResponse{
		Error:   errorType,
		Message: message,
	})
}

func BadRequestError(c *fiber.Ctx, errorType string, message string) error {
	return ErrorResponse(c, fiber.StatusBadRequest, errorType, message)
}

func InternalServerError(c *fiber.Ctx, errorType string, message string) error {
	return ErrorResponse(c, fiber.StatusInternalServerError, errorType, message)
}

func UnauthorizedError(c *fiber.Ctx, errorType string, message string) error {
	return ErrorResponse(c, fiber.StatusUnauthorized, errorType, message)
}

func NotFoundError(c *fiber.Ctx, errorType string, message string) error {
	return ErrorResponse(c, fiber.StatusNotFound, errorType, message)
}

func ForbiddenError(c *fiber.Ctx, errorType string, message string) error {
	return ErrorResponse(c, fiber.StatusForbidden, errorType, message)
}
