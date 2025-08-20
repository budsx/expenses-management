package handler

import (
	"net/http"

	"github.com/budsx/expenses-management/model"
	"github.com/budsx/expenses-management/service"
	"github.com/gofiber/fiber/v2"
)

type AuthHandler struct {
	service *service.ExpensesManagementService
}

func NewAuthHandler(service *service.ExpensesManagementService) *AuthHandler {
	return &AuthHandler{service: service}
}

func (h *AuthHandler) Login(c *fiber.Ctx) error {
	var req model.LoginRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error":   "Invalid request body",
			"message": err.Error(),
		})
	}

	if req.Email == "" || req.Password == "" {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error":   "Missing required fields",
			"message": "Email and password are required",
		})
	}

	loginResponse, err := h.service.AuthenticateUser(c.Context(), req.Email, req.Password)
	if err != nil {
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{
			"error":   "Authentication failed",
			"message": err.Error(),
		})
	}

	return c.Status(http.StatusOK).JSON(model.Response{
		Message: "Success",
		Data:    loginResponse,
	})
}
