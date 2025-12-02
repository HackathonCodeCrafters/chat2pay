package api

import (
	"chat2pay/bootstrap"
	"chat2pay/config/yaml"
	"chat2pay/internal/api/handlers"
	"chat2pay/internal/api/routes"
	"chat2pay/internal/middlewares/jwt"
	"chat2pay/internal/repositories"
	"chat2pay/internal/service"
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

	// Repositories
	merchantRepo := repositories.NewMerchantRepo(db)
	merchantUserRepo := repositories.NewMerchantUserRepo(db)
	productRepo := repositories.NewProductRepo(db)
	customerRepo := repositories.NewCustomerRepo(db)
	orderRepo := repositories.NewOrderRepo(db)

	// Middlewares
	authMdwr := jwt.NewAuthMiddleware(cfg)

	// Services
	merchantService := service.NewMerchantService(merchantRepo, cfg)
	merchantAuthService := service.NewMerchantAuthService(merchantUserRepo, merchantRepo, authMdwr, cfg)
	customerAuthService := service.NewCustomerAuthService(customerRepo, authMdwr, cfg)
	productService := service.NewProductService(productRepo, merchantRepo, cfg)
	customerService := service.NewCustomerService(customerRepo, cfg)
	orderService := service.NewOrderService(orderRepo, productRepo, customerRepo, merchantRepo, cfg)

	// Handlers
	merchantAuthHandler := handlers.NewMerchantAuthHandler(merchantAuthService)
	customerAuthHandler := handlers.NewCustomerAuthHandler(customerAuthService)
	merchantHandler := handlers.NewMerchantHandler(merchantService)
	productHandler := handlers.NewProductHandler(productService)
	customerHandler := handlers.NewCustomerHandler(customerService)
	orderHandler := handlers.NewOrderHandler(orderService)

	// API Group
	api := router.Group("/api")

	router.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"message": "Chat2Pay Backend API",
			"version": "1.0.0",
			"status":  "running",
		})
	})

	// Routes
	routes.AuthRouter(api, merchantAuthHandler, customerAuthHandler)
	routes.MerchantRouter(api, merchantHandler, authMdwr)
	routes.ProductRouter(api, productHandler, authMdwr)
	routes.CustomerRouter(api, customerHandler, authMdwr)
	routes.OrderRouter(api, orderHandler, authMdwr)

	return router
}
