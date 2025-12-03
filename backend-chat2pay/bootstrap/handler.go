package bootstrap

import (
	"chat2pay/internal/api/handlers"
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
	}
}
