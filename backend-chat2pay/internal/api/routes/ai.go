package routes

import (
	"chat2pay/config/yaml"
	"chat2pay/internal/api/handlers"
	"chat2pay/internal/service"
	"github.com/gofiber/fiber/v2"
)

func AiRoutes(app fiber.Router, cfg *yaml.Config, service service.AiService) {
	app.Post("/prompt", handlers.Prompt(cfg, service))
}
