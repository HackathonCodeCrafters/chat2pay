package routes

import (
	"chat2pay/internal/api/handlers"
	"chat2pay/internal/middlewares/jwt"
	"github.com/gofiber/fiber/v2"
)

func CustomerRouter(router fiber.Router, handler *handlers.CustomerHandler, authMdwr jwt.AuthMiddleware) {
	customers := router.Group("/customers")

	// Protected routes - merchant can view all, customer can view self
	customers.Get("/", jwt.JWTProtected(authMdwr), jwt.RequireRole("merchant"), handler.GetAll)
	customers.Get("/:id", jwt.JWTProtected(authMdwr), handler.GetById)

	// Protected routes - authenticated users
	customers.Put("/:id", jwt.JWTProtected(authMdwr), handler.Update)
	customers.Delete("/:id", jwt.JWTProtected(authMdwr), handler.Delete)
}
