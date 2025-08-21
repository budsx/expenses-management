package handler

import (
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
		return BadRequestError(c, "Invalid request body", err.Error())
	}

	if req.Email == "" || req.Password == "" {
		return BadRequestError(c, "Missing required fields", "Email and password are required")
	}

	resp, err := h.service.AuthenticateUser(c.Context(), req.Email, req.Password)
	if err != nil {
		return UnauthorizedError(c, "Authentication failed", err.Error())
	}

	return SuccessResponse(c, "success", resp)
}
