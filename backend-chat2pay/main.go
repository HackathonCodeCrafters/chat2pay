package main

import (
	"chat2pay/bootstrap"
	command "chat2pay/cmd"
	"chat2pay/config/yaml"
	"chat2pay/internal/api"
	"context"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/urfave/cli/v3"
	"os"

	"log"
)

// @title Chat2Pay API
// @version 1.0
// @description Chat2Pay Backend API - E-commerce platform dengan fitur chat-to-pay menggunakan AI/LLM
// @termsOfService http://swagger.io/terms/

// @contact.name Chat2Pay Team
// @contact.email support@chat2pay.com

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost:9005
// @BasePath /api

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Masukkan token dengan format: Bearer {token}

func main() {

	ctn, err := bootstrap.NewContainer()
	if err != nil {
		panic(err)
	}

	cmd := cli.Command{}

	cmd.Commands = append(cmd.Commands,
		&cli.Command{
			Name:  "http",
			Usage: "Run Chat2Pay http",
			Action: func(context.Context, *cli.Command) error {
				config := ctn.Get(bootstrap.ConfigDefName).(*yaml.Config)

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

				app = api.NewRouter(ctn)

				log.Fatal(app.Listen(fmt.Sprintf(`:%s`, config.App.Port)))
				return nil
			},
		},
	)

	cmd.Commands = append(cmd.Commands, command.Migration(&ctn)...)

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}
