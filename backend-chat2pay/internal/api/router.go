package api

import (
	"chat2pay/bootstrap"
	"chat2pay/internal/api/handlers"
	"chat2pay/internal/api/routes"
	"github.com/gofiber/fiber/v2"
	"github.com/sarulabs/di/v2"
)

func NewRouter(ctn di.Container) *fiber.App {
	router := fiber.New()

	// Socket
	handlers.NewSocketEvent(router, ctn)

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
	//api.AuthRouter(api, merchantAuthHandler, customerAuthHandler)
	//routes.MerchantRouter(api, merchantHandler, authMdwr)
	routes.ProductRouter(api, ctn.Get(bootstrap.ProductHandlerName).(*handlers.ProductHandler))
	//routes.CustomerRouter(api, customerHandler, authMdwr)
	//routes.OrderRouter(api, orderHandler, authMdwr)

	return router
}
