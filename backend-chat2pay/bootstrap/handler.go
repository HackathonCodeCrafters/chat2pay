package bootstrap

import (
	"chat2pay/config/yaml"
	"chat2pay/internal/api/handlers"
	"chat2pay/internal/repositories"
	"chat2pay/internal/service"
	"github.com/sarulabs/di/v2"
)

func LoadHandler() *[]di.Def {
	return &[]di.Def{
		{
			Name: ProductHandlerName,
			Build: func(ctn di.Container) (interface{}, error) {
				productService := ctn.Get(ProductServiceName).(service.ProductService)
				return handlers.NewProductHandler(productService), nil
			},
		},
		{
			Name: CustomerHandlerName,
			Build: func(ctn di.Container) (interface{}, error) {
				productService := ctn.Get(ProductServiceName).(service.ProductService)
				return handlers.NewProductHandler(productService), nil
			},
		},
		{
			Name: CustomerAuthHandlerName,
			Build: func(ctn di.Container) (interface{}, error) {
				customerAuthService := ctn.Get(AuthCustomerService).(service.CustomerAuthService)
				return handlers.NewCustomerAuthHandler(customerAuthService), nil
			},
		},
		{
			Name: MerchantAuthHandlerName,
			Build: func(ctn di.Container) (interface{}, error) {
				merchantAuthService := ctn.Get(AuthMerchantService).(service.MerchantAuthService)
				return handlers.NewMerchantAuthHandler(merchantAuthService), nil
			},
		},
		{
			Name: ShippingHandlerName,
			Build: func(ctn di.Container) (interface{}, error) {
				cfg := ctn.Get(ConfigDefName).(*yaml.Config)
				return handlers.NewShippingHandler(cfg.RajaOngkir.APIKey), nil
			},
		},
		{
			Name: OrderHandlerName,
			Build: func(ctn di.Container) (interface{}, error) {
				orderService := ctn.Get(OrderServiceName).(service.OrderService)
				return handlers.NewOrderHandler(orderService), nil
			},
		},
		{
			Name: ChatHandlerName,
			Build: func(ctn di.Container) (interface{}, error) {
				chatRepo := ctn.Get(ChatMessageRepositoryName).(repositories.ChatMessageRepository)
				return handlers.NewChatHandler(chatRepo), nil
			},
		},
	}
}
