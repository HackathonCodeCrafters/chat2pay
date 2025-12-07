package routes

import (
	"chat2pay/internal/api/handlers"
	"github.com/gofiber/fiber/v2"
)

func AuthRouter(
	router fiber.Router,
	merchantAuthHandler *handlers.MerchantAuthHandler,
	customerAuthHandler *handlers.CustomerAuthHandler,
) {
	auth := router.Group("/auth")

	// Merchant auth
	if merchantAuthHandler != nil {
		merchant := auth.Group("/merchant")
		merchant.Post("/register", merchantAuthHandler.Register)
		merchant.Post("/login", merchantAuthHandler.Login)
	}

	// Customer auth
	if customerAuthHandler != nil {
		customer := auth.Group("/customer")
		customer.Post("/register", customerAuthHandler.Register)
		customer.Post("/login", customerAuthHandler.Login)
	}
}
