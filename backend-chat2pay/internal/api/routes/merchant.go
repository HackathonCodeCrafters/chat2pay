package routes

import (
	"chat2pay/internal/api/handlers"
	"chat2pay/internal/middlewares/jwt"
	"github.com/gofiber/fiber/v2"
)

func MerchantRouter(router fiber.Router, handler *handlers.MerchantHandler, authMdwr jwt.AuthMiddleware) {
	merchants := router.Group("/merchants")

	// Public routes
	merchants.Get("/", handler.GetAll)
	merchants.Get("/:id", handler.GetById)

	// Protected routes - merchant only
	merchants.Post("/", jwt.JWTProtected(authMdwr), jwt.RequireRole("merchant"), handler.Create)
	merchants.Put("/:id", jwt.JWTProtected(authMdwr), jwt.RequireRole("merchant"), handler.Update)
	merchants.Delete("/:id", jwt.JWTProtected(authMdwr), jwt.RequireRole("merchant"), handler.Delete)
}
