package handler

import (
	"github.com/budsx/expenses-management/service"
)

type AuthHandler struct {
	service *service.ExpensesManagementService
}

func NewAuthHandler(service *service.ExpensesManagementService) *AuthHandler {
	return &AuthHandler{service: service}
}
