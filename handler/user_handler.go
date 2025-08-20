package handler

import (
	"github.com/budsx/expenses-management/service"
	"github.com/gofiber/fiber/v2"
)

func NewUserHandler(service *service.ExpensesManagementService) *UserHandler {
	return &UserHandler{service: service}
}

type UserHandler struct {
	service *service.ExpensesManagementService
}

func (h *UserHandler) GetUser(c *fiber.Ctx) error {
	user, err := h.service.GetUser(c.Context(), "1")
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"message": "User fetched successfully", "user": user})
}
