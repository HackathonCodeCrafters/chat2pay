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
				fmt.Println("kesini blog goblg")
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
