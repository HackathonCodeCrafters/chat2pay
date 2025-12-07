package routes

import (
	"chat2pay/internal/api/handlers"
	"chat2pay/internal/api/middleware"
	"github.com/gofiber/fiber/v2"
)

func OrderRouter(router fiber.Router, handler *handlers.OrderHandler, jwtSecret string) {
	orders := router.Group("/orders")

	customerAuth := middleware.CustomerAuthMiddleware(jwtSecret)
	merchantAuth := middleware.MerchantAuthMiddleware(jwtSecret)

	// Static routes MUST come before parameterized routes
	// Customer routes
	orders.Post("/", customerAuth, handler.CreateOrder)
	orders.Get("/customer", customerAuth, handler.GetCustomerOrders)

	// Merchant routes
	orders.Get("/merchant", merchantAuth, handler.GetMerchantOrders)
	orders.Get("/merchant/:id", merchantAuth, handler.GetOrder)
	orders.Patch("/:id/status", merchantAuth, handler.UpdateOrderStatus)

	// Parameterized routes LAST
	orders.Get("/:id", customerAuth, handler.GetOrder)
}
