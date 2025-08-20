package http

import (
	"github.com/budsx/expenses-management/handler"
	"github.com/budsx/expenses-management/service"
	"github.com/gofiber/fiber/v2"
)

type ExpensesManagementServer struct {
	app             *fiber.App
	userHandler     *handler.UserHandler
	expensesHandler *handler.ExpensesManagementHandler
}

func NewExpensesManagementServer(service *service.ExpensesManagementService, userHandler *handler.UserHandler, expensesHandler *handler.ExpensesManagementHandler) *ExpensesManagementServer {
	app := fiber.New()

	// Setup Routes
	app.Get("/user", userHandler.GetUser)

	return &ExpensesManagementServer{
		app:             app,
		userHandler:     userHandler,
		expensesHandler: expensesHandler,
	}
}

func (s *ExpensesManagementServer) Run(port string) error {
	return s.app.Listen(port)
}
