package routes

import (
	"chat2pay/internal/api/handlers"
	"github.com/gofiber/fiber/v2"
)

func ShippingRouter(router fiber.Router, handler *handlers.ShippingHandler) {
	shipping := router.Group("/shipping")

	shipping.Get("/provinces", handler.GetProvinces)
	shipping.Get("/cities", handler.GetCities)
	shipping.Get("/cost", handler.GetCost)
	shipping.Get("/track", handler.TrackWaybill)
}
