package routes

import (
	"chat2pay/internal/api/handlers"
	"github.com/gofiber/fiber/v2"
)

func ProductRouter(router fiber.Router, handler *handlers.ProductHandler) {
	products := router.Group("/products")

	// Public routes
	products.Get("/", handler.GetAll)
	products.Get("/:id", handler.GetById)
	products.Post("/", handler.Create)
	products.Post("/multiple", handler.CreateMultiple)
	products.Post("/ask", handler.AskProduct)

	//// Protected routes - merchant only
	//products.Post("/", jwt.JWTProtected(authMdwr), jwt.RequireRole("merchant"), handler.Create)
	//products.Put("/:id", jwt.JWTProtected(authMdwr), jwt.RequireRole("merchant"), handler.Update)
	//products.Delete("/:id", jwt.JWTProtected(authMdwr), jwt.RequireRole("merchant"), handler.Delete)
}
