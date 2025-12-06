package bootstrap

import (
	"chat2pay/config/yaml"
	"chat2pay/internal/pkg/llm/mistral"
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
				//geminiModel := ctn.Get(GeminiLLMName).(*gemini.GeminiLLM)
				//openAIModel := ctn.Get(OpenAILLMName).(*openai.OpenAI)
				mistralModel := ctn.Get(MistralLLMName).(*mistral.MistralLLM)

				config := ctn.Get(ConfigDefName).(*yaml.Config)
				return service.NewProductService(productRepo, merchantRepo, mistralModel, config), nil
			},
		},
	}
}
