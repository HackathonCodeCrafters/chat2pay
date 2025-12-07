package api

import (
	"chat2pay/bootstrap"
	"chat2pay/config/yaml"
	_ "chat2pay/docs"
	"chat2pay/internal/api/handlers"
	"chat2pay/internal/api/routes"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/swagger"
	"github.com/sarulabs/di/v2"
)

func NewRouter(ctn di.Container) *fiber.App {
	router := fiber.New()

	// CORS middleware
	router.Use(cors.New(cors.Config{
		AllowOrigins:     "*",
		AllowMethods:     "GET,POST,PUT,PATCH,DELETE,OPTIONS",
		AllowHeaders:     "Origin,Content-Type,Accept,Authorization",
		AllowCredentials: false,
	}))

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

	// Routes
	config := ctn.Get(bootstrap.ConfigDefName).(*yaml.Config)

	routes.AuthRouter(
		api,
		ctn.Get(bootstrap.MerchantAuthHandlerName).(*handlers.MerchantAuthHandler),
		ctn.Get(bootstrap.CustomerAuthHandlerName).(*handlers.CustomerAuthHandler),
	)
	routes.ProductRouter(api, ctn.Get(bootstrap.ProductHandlerName).(*handlers.ProductHandler))
	routes.ShippingRouter(api, ctn.Get(bootstrap.ShippingHandlerName).(*handlers.ShippingHandler))
	routes.OrderRouter(api, ctn.Get(bootstrap.OrderHandlerName).(*handlers.OrderHandler), config.JWT.Key)
	routes.ChatRouter(api, ctn.Get(bootstrap.ChatHandlerName).(*handlers.ChatHandler), config.JWT.Key)

	return router
}
