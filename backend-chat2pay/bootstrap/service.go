package bootstrap

import (
	"chat2pay/config/yaml"
	"chat2pay/internal/repositories"
	"chat2pay/internal/service"
	"github.com/sarulabs/di/v2"
)

func LoadService() *[]di.Def {
	return &[]di.Def{
		{
			Name: ProductServiceName,
			Build: func(ctn di.Container) (interface{}, error) {
				productRepo := ctn.Get(ProductRepositoryName).(repositories.ProductRepository)
				merchantRepo := ctn.Get(MerchantRepositoryName).(repositories.MerchantRepository)
				config := ctn.Get(ConfigDefName).(*yaml.Config)
				return service.NewProductService(productRepo, merchantRepo, config), nil
			},
		},
	}
}
