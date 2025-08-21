package handler

import (
	"strings"

	"github.com/budsx/expenses-management/util"
	"github.com/gofiber/fiber/v2"
)

func AuthMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return UnauthorizedError(c, "Authorization header required", "Authorization header required")
		}

		if !strings.HasPrefix(authHeader, "Bearer ") {
			return UnauthorizedError(c, "Authorization header must start with Bearer", "Authorization header must start with Bearer")
		}

		token := strings.TrimPrefix(authHeader, "Bearer ")
		if token == "" {
			return UnauthorizedError(c, "Token not provided", "Token not provided")
		}

		claims, err := util.ValidateJWT(token)
		if err != nil {
			return UnauthorizedError(c, "Invalid token", err.Error())
		}

		c.Locals("user_id", claims.UserID)
		c.Locals("user_email", claims.Email)
		c.Locals("user_role", claims.Role)

		return c.Next()
	}
}
