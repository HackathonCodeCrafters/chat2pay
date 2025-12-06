package api

import (
	"chat2pay/bootstrap"
	_ "chat2pay/docs"
	"chat2pay/internal/api/handlers"
	"chat2pay/internal/api/routes"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"
	"github.com/sarulabs/di/v2"
)

func NewRouter(ctn di.Container) *fiber.App {
	router := fiber.New()

	// Swagger UI
	router.Get("/swagger/*", swagger.HandlerDefault)

	// API Group
	api := router.Group("/api")

	api.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"message": "Chat2Pay Backend API",
			"version": "1.0.0",
			"status":  "running",
		})
	})

	// Payment routes
	paymentHandler := ctn.Get(bootstrap.PaymentHandlerName).(*handlers.PaymentHandler)
	webhookHandler := ctn.Get(bootstrap.WebhookHandlerName).(*handlers.WebhookHandler)
	routes.PaymentRoutes(api, paymentHandler, webhookHandler)

	return router
}
