package routes

import (
	"chat2pay/config/yaml"
	"chat2pay/internal/api/handlers"
	"chat2pay/internal/middlewares/jwt"
	"chat2pay/internal/service"
	"github.com/gofiber/fiber/v2"
)

func AuthRouter(app fiber.Router, cfg *yaml.Config, authMiddleware jwt.AuthMiddleware, service service.AuthService) {
	app.Post("/login", handlers.Login(cfg, service))
	app.Get("/profile", authMiddleware.ValidateToken(), handlers.GetProfile(cfg, service))
}
