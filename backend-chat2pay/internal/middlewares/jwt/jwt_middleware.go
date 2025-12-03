package jwt

import (
	"chat2pay/internal/api/presenter"
	"errors"
	"strings"

	"github.com/gofiber/fiber/v2"
)

func JWTProtected(authMdwr AuthMiddleware) fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(401).JSON(presenter.ErrorResponse(errors.New("missing authorization header")))
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			return c.Status(401).JSON(presenter.ErrorResponse(errors.New("invalid authorization header format")))
		}

		tokenString := parts[1]
		claims, err := authMdwr.ValidateToken(tokenString)
		if err != nil {
			return c.Status(401).JSON(presenter.ErrorResponse(errors.New("invalid or expired token")))
		}

		c.Locals("user_id", claims.UserID)
		c.Locals("email", claims.Email)
		c.Locals("role", claims.Role)

		return c.Next()
	}
}

func RequireRole(roles ...string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		userRole := c.Locals("role").(string)

		for _, role := range roles {
			if userRole == role {
				return c.Next()
			}
		}

		return c.Status(403).JSON(presenter.ErrorResponse(errors.New("insufficient permissions")))
	}
}
