package bootstrap

import (
	"chat2pay/config/yaml"
	"chat2pay/internal/pkg/llm/gemini"
	"chat2pay/internal/pkg/llm/mistral"
	"chat2pay/internal/pkg/llm/openai"
	"chat2pay/internal/pkg/rajaongkir"
	"github.com/sarulabs/di/v2"
)

func LoadPackage() *[]di.Def {
	return &[]di.Def{
		{
			Name: GeminiLLMName,
			Build: func(ctn di.Container) (interface{}, error) {
				config := ctn.Get(ConfigDefName).(*yaml.Config)
				return gemini.NewGeminiLLM(config.Gemini.APIKey), nil
			},
		},
		{
			Name: OpenAILLMName,
			Build: func(ctn di.Container) (interface{}, error) {
				config := ctn.Get(ConfigDefName).(*yaml.Config)
				return openai.NewOpenAI(config.OpenAI.APIKey), nil
			},
		},
		{
			Name: MistralLLMName,
			Build: func(ctn di.Container) (interface{}, error) {
				config := ctn.Get(ConfigDefName).(*yaml.Config)
				return mistral.NewMistralLLM(config.Mistral.APIKey), nil
			},
		},
		{
			Name: RajaOngkirName,
			Build: func(ctn di.Container) (interface{}, error) {
				config := ctn.Get(ConfigDefName).(*yaml.Config)
				return rajaongkir.NewRajaOngkir(config.RajaOngkir.APIKey), nil
			},
		},
	}
}
