package bootstrap

import (
	"chat2pay/config/yaml"
	"chat2pay/internal/middlewares/jwt"
	"chat2pay/internal/pkg/llm/mistral"
	"chat2pay/internal/repositories"
	"chat2pay/internal/service"
	"github.com/sarulabs/di/v2"
)

func LoadService() *[]di.Def {
	return &[]di.Def{
		{
			Name: AuthMiddlewareName,
			Build: func(ctn di.Container) (interface{}, error) {
				config := ctn.Get(ConfigDefName).(*yaml.Config)
				return jwt.NewAuthMiddleware(config), nil
			},
		},
		{
			Name: ProductServiceName,
			Build: func(ctn di.Container) (interface{}, error) {
				productRepo := ctn.Get(ProductRepositoryName).(repositories.ProductRepository)
				merchantRepo := ctn.Get(MerchantRepositoryName).(repositories.MerchantRepository)
				//geminiModel := ctn.Get(GeminiLLMName).(*gemini.GeminiLLM)
				//openAIModel := ctn.Get(OpenAILLMName).(*openai.OpenAI)
				mistralModel := ctn.Get(MistralLLMName).(*mistral.MistralLLM)

				config := ctn.Get(ConfigDefName).(*yaml.Config)
				return service.NewProductService(productRepo, merchantRepo, mistralModel, config), nil
			},
		},
		{
			Name: AuthCustomerService,
			Build: func(ctn di.Container) (interface{}, error) {
				customerRepo := ctn.Get(CustomerRepositoryName).(repositories.CustomerRepository)
				authMdwr := ctn.Get(AuthMiddlewareName).(jwt.AuthMiddleware)
				config := ctn.Get(ConfigDefName).(*yaml.Config)
				return service.NewCustomerAuthService(customerRepo, authMdwr, config), nil
			},
		},
		{
			Name: AuthMerchantService,
			Build: func(ctn di.Container) (interface{}, error) {
				merchantUserRepo := ctn.Get(MerchantUserRepositoryName).(repositories.MerchantUserRepository)
				merchantRepo := ctn.Get(MerchantRepositoryName).(repositories.MerchantRepository)
				authMdwr := ctn.Get(AuthMiddlewareName).(jwt.AuthMiddleware)
				config := ctn.Get(ConfigDefName).(*yaml.Config)
				return service.NewMerchantAuthService(merchantUserRepo, merchantRepo, authMdwr, config), nil
			},
		},
		{
			Name: OrderServiceName,
			Build: func(ctn di.Container) (interface{}, error) {
				orderRepo := ctn.Get(OrderRepositoryName).(repositories.OrderRepository)
				productRepo := ctn.Get(ProductRepositoryName).(repositories.ProductRepository)
				config := ctn.Get(ConfigDefName).(*yaml.Config)
				return service.NewOrderService(config, orderRepo, productRepo), nil
			},
		},
	}
}
