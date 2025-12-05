package bootstrap

import (
	"chat2pay/config/yaml"
	"chat2pay/internal/pkg/llm/gemini"
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
	}
}
