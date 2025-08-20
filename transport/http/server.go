package http

import (
	"github.com/budsx/expenses-management/handler"
	"github.com/budsx/expenses-management/service"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
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

	// Middleware
	app.Use(logger.New())
	app.Use(cors.New())

	// Setup Routes
	// Health check endpoint
	app.Get("/api/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":  "healthy",
			"service": "expenses-management",
		})
	})

	// Direct auth routes (for backward compatibility)
	app.Post("/auth/login", authHandler.Login)

	// Auth routes
	auth := app.Group("/auth")
	auth.Post("/login", authHandler.Login)

	// User routes
	users := app.Group("/users")
	users.Get("/:id", userHandler.GetUser)

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
