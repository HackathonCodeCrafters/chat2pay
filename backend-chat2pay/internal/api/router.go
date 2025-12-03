package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/sarulabs/di/v2"
)

func NewRouter(container di.Container) *fiber.App {
	router := fiber.New()

	// API Group
	api := router.Group("/api")

	api.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"message": "Chat2Pay Backend API",
			"version": "1.0.0",
			"status":  "running",
		})
	})

	// Routes
	//routes.AuthRouter(api, merchantAuthHandler, customerAuthHandler)
	//routes.MerchantRouter(api, merchantHandler, authMdwr)
	//routes.ProductRouter(api, productHandler, authMdwr)
	//routes.CustomerRouter(api, customerHandler, authMdwr)
	//routes.OrderRouter(api, orderHandler, authMdwr)

	return router
}
