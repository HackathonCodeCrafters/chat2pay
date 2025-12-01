package api

import (
	"chat2pay/bootstrap"
	"chat2pay/config/yaml"
	"chat2pay/migrations"
	"fmt"
	"github.com/gofiber/fiber/v2"
)

func NewRouter(cfg *yaml.Config) *fiber.App {
	router := fiber.New()

	db, err := bootstrap.DatabaseConnection(cfg)

	if err != nil {
		panic(fmt.Sprintf(`db connection error got : %v`, err))
	}

	fmt.Println("Database connection success!")

	migrations.AutoMigration(db)

	//repositories
	var (
	//userRepo = repositories.NewUserRepo(db)
	)

	//middlewares
	var (
	//authMidleware = jwt.NewAuthMiddleware(userRepo, cfg)
	)

	//service
	var (
	//authService = service.NewAuthService(authMidleware, userRepo, cfg)
	//aiService = service.NewAiService(cfg, models)
	)

	//group
	//api := router.Group("/api")

	router.Get("/", func(c *fiber.Ctx) error {
		return c.Send([]byte("uwow"))
	})

	//routes.AuthRouter(api, cfg, authMidleware, authService)
	//routes.AiRoutes(api, cfg, aiService)

	//routes.HealthRouter(api)

	return router
}
