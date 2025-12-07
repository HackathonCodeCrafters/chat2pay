package routes

import (
	"chat2pay/internal/api/handlers"
	"chat2pay/internal/api/middleware"
	"github.com/gofiber/fiber/v2"
)

func ChatRouter(router fiber.Router, handler *handlers.ChatHandler, jwtSecret string) {
	chat := router.Group("/chat")
	
	customerAuth := middleware.CustomerAuthMiddleware(jwtSecret)
	
	chat.Get("/history", customerAuth, handler.GetHistory)
	chat.Post("/messages", customerAuth, handler.SaveMessage)
	chat.Delete("/history", customerAuth, handler.ClearHistory)
}
