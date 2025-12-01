package main

import (
	"chat2pay/config/yaml"
	"chat2pay/internal/api"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"

	"log"
)

func main() {

	cfg, err := yaml.NewConfig()
	if err != nil {
		log.Fatal(fmt.Sprintf(`read cfg yaml got error : %v`, err))
	}

	app := fiber.New()
	app.Use(cors.New())

	// Or extend your config for customization
	app.Use(cors.New(cors.Config{
		AllowHeaders:     "Origin,Content-Type,Accept,Content-Length,Accept-Language,Accept-Encoding,Connection,Access-Control-Allow-Origin",
		AllowOrigins:     "*",
		AllowCredentials: false,
		AllowMethods:     "GET,POST",
		MaxAge:           3,
	}))

	app = api.NewRouter(cfg)

	log.Fatal(app.Listen(fmt.Sprintf(`:%s`, cfg.App.Port)))
}
