package http

import (
	"github.com/budsx/expenses-management/handler"
	"github.com/budsx/expenses-management/middleware"
	"github.com/budsx/expenses-management/service"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/google/uuid"
)

type ExpensesManagementServer struct {
	app             *fiber.App
	userHandler     *handler.UserHandler
	expensesHandler *handler.ExpensesManagementHandler
	authHandler     *handler.AuthHandler
}

func NewExpensesManagementServer(service *service.ExpensesManagementService, userHandler *handler.UserHandler, expensesHandler *handler.ExpensesManagementHandler, authHandler *handler.AuthHandler) *ExpensesManagementServer {
	app := fiber.New(fiber.Config{
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			code := fiber.StatusInternalServerError
			if e, ok := err.(*fiber.Error); ok {
				code = e.Code
			}
			return c.Status(code).JSON(fiber.Map{
				"error":   "Internal Server Error",
				"message": err.Error(),
			})
		},
	})

	app.Use(cors.New())

	app.Use(func(c *fiber.Ctx) error {
		c.Locals("request_id", uuid.New().String())
		return c.Next()
	})

	app.Get("/api/health", func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"status":  "healthy",
			"service": "expenses-management",
		})
	})

	api := app.Group("/api")

	api.Post("/auth/login", authHandler.Login)

	api.Use(middleware.AuthMiddleware())
	api.Get("/users/:id", userHandler.GetUser)

	expenses := api.Group("/expenses")
	expenses.Use(middleware.AuthMiddleware())
	expenses.Post("/", expensesHandler.CreateExpense)
	expenses.Get("/", expensesHandler.GetExpenses)
	expenses.Get("/:id", expensesHandler.GetExpenseByID)
	expenses.Put("/:id/approve", expensesHandler.ApproveExpense)
	expenses.Put("/:id/reject", expensesHandler.RejectExpense)

	return &ExpensesManagementServer{
		app:             app,
		userHandler:     userHandler,
		expensesHandler: expensesHandler,
		authHandler:     authHandler,
	}
}

func (s *ExpensesManagementServer) Run(port string) error {
	return s.app.Listen(port)
}
