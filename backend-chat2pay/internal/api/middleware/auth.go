package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"strings"
)

func CustomerAuthMiddleware(jwtSecret string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(401).JSON(fiber.Map{
				"status": false,
				"error":  "Missing authorization header",
			})
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader {
			return c.Status(401).JSON(fiber.Map{
				"status": false,
				"error":  "Invalid authorization format",
			})
		}

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return []byte(jwtSecret), nil
		})

		if err != nil || !token.Valid {
			return c.Status(401).JSON(fiber.Map{
				"status": false,
				"error":  "Invalid or expired token",
			})
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			return c.Status(401).JSON(fiber.Map{
				"status": false,
				"error":  "Invalid token claims",
			})
		}

		// Check role and set customer info in context
		role, _ := claims["role"].(string)
		if role != "customer" {
			return c.Status(401).JSON(fiber.Map{
				"status": false,
				"error":  "Invalid token: customer token required",
			})
		}

		if userID, ok := claims["user_id"].(string); ok {
			c.Locals("customer_id", userID)
		}
		if email, ok := claims["email"].(string); ok {
			c.Locals("email", email)
		}

		return c.Next()
	}
}

func MerchantAuthMiddleware(jwtSecret string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(401).JSON(fiber.Map{
				"status": false,
				"error":  "Missing authorization header",
			})
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader {
			return c.Status(401).JSON(fiber.Map{
				"status": false,
				"error":  "Invalid authorization format",
			})
		}

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return []byte(jwtSecret), nil
		})

		if err != nil || !token.Valid {
			return c.Status(401).JSON(fiber.Map{
				"status": false,
				"error":  "Invalid or expired token",
			})
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			return c.Status(401).JSON(fiber.Map{
				"status": false,
				"error":  "Invalid token claims",
			})
		}

		// Check role and set merchant info in context
		role, _ := claims["role"].(string)
		if role != "merchant" {
			return c.Status(401).JSON(fiber.Map{
				"status": false,
				"error":  "Invalid token: merchant token required",
			})
		}

		if merchantID, ok := claims["merchant_id"].(string); ok {
			c.Locals("merchant_id", merchantID)
		}
		if userID, ok := claims["user_id"].(string); ok {
			c.Locals("user_id", userID)
		}
		if email, ok := claims["email"].(string); ok {
			c.Locals("email", email)
		}

		return c.Next()
	}
}
