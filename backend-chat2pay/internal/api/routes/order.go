package routes

import (
	"chat2pay/internal/api/handlers"
	"chat2pay/internal/middlewares/jwt"
	"github.com/gofiber/fiber/v2"
)

func OrderRouter(router fiber.Router, handler *handlers.OrderHandler, authMdwr jwt.AuthMiddleware) {
	orders := router.Group("/orders")

	// All order routes require authentication
	orders.Post("/", jwt.JWTProtected(authMdwr), handler.Create)
	orders.Get("/", jwt.JWTProtected(authMdwr), handler.GetAll)
	orders.Get("/:id", jwt.JWTProtected(authMdwr), handler.GetById)

	// Only merchant can update status and delete
	orders.Patch("/:id/status", jwt.JWTProtected(authMdwr), jwt.RequireRole("merchant"), handler.UpdateStatus)
	orders.Delete("/:id", jwt.JWTProtected(authMdwr), jwt.RequireRole("merchant"), handler.Delete)
}
