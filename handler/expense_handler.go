package handler

import (
	"strconv"

	"github.com/budsx/expenses-management/model"
	"github.com/budsx/expenses-management/service"

	"github.com/gofiber/fiber/v2"
)

type ExpensesManagementHandler struct {
	service *service.ExpensesManagementService
}

func NewExpensesManagementHandler(service *service.ExpensesManagementService) *ExpensesManagementHandler {
	return &ExpensesManagementHandler{service: service}
}

func (h *ExpensesManagementHandler) HealthCheck(c *fiber.Ctx) error {
	err := h.service.HealthCheck(c.Context())
	if err != nil {
		return InternalServerError(c, "Failed to check health", err.Error())
	}
	return c.JSON(fiber.Map{
		"message": "success",
	})
}

func (h *ExpensesManagementHandler) CreateExpense(c *fiber.Ctx) error {
	var req model.CreateExpenseRequest

	if err := c.BodyParser(&req); err != nil {
		return BadRequestError(c, "Invalid request body", err.Error())
	}

	result, err := h.service.CreateExpense(c.Context(), req)
	if err != nil {
		return InternalServerError(c, "Failed to create expense", err.Error())
	}

	return SuccessResponse(c, "success", result)
}

func (h *ExpensesManagementHandler) GetExpenses(c *fiber.Ctx) error {
	var query model.ExpenseListQuery

	if err := c.QueryParser(&query); err != nil {
		return BadRequestError(c, "Invalid query parameters", err.Error())
	}

	if query.Page <= 0 {
		query.Page = 1
	}
	if query.PageSize <= 0 {
		query.PageSize = 10
	}

	result, err := h.service.GetExpenses(c.Context(), query)
	if err != nil {
		return InternalServerError(c, "Failed to get expenses", err.Error())
	}

	return SuccessResponse(c, "success", result)
}

func (h *ExpensesManagementHandler) GetExpenseByID(c *fiber.Ctx) error {
	expenseIDStr := c.Params("id")
	expenseID, err := strconv.ParseInt(expenseIDStr, 10, 64)
	if err != nil {
		return BadRequestError(c, "Invalid expense ID", "Expense ID must be a valid number")
	}

	result, err := h.service.GetExpenseByID(c.Context(), expenseID)
	if err != nil {
		return InternalServerError(c, "Failed to get expense", err.Error())
	}

	return SuccessResponse(c, "success", result)
}

func (h *ExpensesManagementHandler) ApproveExpense(c *fiber.Ctx) error {
	expenseIDStr := c.Params("id")
	expenseID, err := strconv.ParseInt(expenseIDStr, 10, 64)
	if err != nil {
		return BadRequestError(c, "Invalid expense ID", "Expense ID must be a valid number")
	}

	req := model.ApprovalRequest{}
	if err := c.BodyParser(&req); err != nil {
		return BadRequestError(c, "Invalid request body", err.Error())
	}
	req.ExpenseID = expenseID

	result, err := h.service.ApproveExpense(c.Context(), req)
	if err != nil {
		return InternalServerError(c, "Failed to approve expense", err.Error())
	}

	return SuccessResponse(c, "success", result)
}

func (h *ExpensesManagementHandler) RejectExpense(c *fiber.Ctx) error {
	expenseIDStr := c.Params("id")
	expenseID, err := strconv.ParseInt(expenseIDStr, 10, 64)
	if err != nil {
		return BadRequestError(c, "Invalid expense ID", "Expense ID must be a valid number")
	}

	req := model.ApprovalRequest{}
	if err := c.BodyParser(&req); err != nil {
		return BadRequestError(c, "Invalid request body", err.Error())
	}
	req.ExpenseID = expenseID

	result, err := h.service.RejectExpense(c.Context(), req)
	if err != nil {
		return InternalServerError(c, "Failed to reject expense", err.Error())
	}

	return SuccessResponse(c, "success", result)
}
