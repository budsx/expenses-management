package handler

import "github.com/budsx/expenses-management/service"

type ExpensesManagementHandler struct {
	service *service.ExpensesManagementService
}

func NewExpensesManagementHandler(service *service.ExpensesManagementService) *ExpensesManagementHandler {
	return &ExpensesManagementHandler{service: service}
}