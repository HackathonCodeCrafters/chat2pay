package bootstrap

import (
	"chat2pay/config/yaml"
	"chat2pay/internal/pkg/llm/gemini"
	"chat2pay/internal/pkg/llm/mistral"
	"chat2pay/internal/pkg/llm/openai"
	"chat2pay/internal/pkg/rajaongkir"
	"chat2pay/internal/pkg/llm"
	"chat2pay/internal/pkg/redis"
	"github.com/sarulabs/di/v2"
)

func LoadPackage() *[]di.Def {
	return &[]di.Def{
		{
			Name: LLMPackageName,
			Build: func(ctn di.Container) (interface{}, error) {
				config := ctn.Get(ConfigDefName).(*yaml.Config)
				redisClient := ctn.Get(RedisAdapter).(redis.RedisClient)
				return llm.NewLLM(config, redisClient), nil
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
