package bootstrap

import (
	"chat2pay/config/yaml"
	"chat2pay/internal/service"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/sarulabs/di/v2"
	"time"
)

func loadService(builder *di.Builder, config *yaml.Config) {
	builder.Add([]di.Def{
		{
			Name: ProductServiceName,
			Build: func(ctn di.Container) (interface{}, error) {
				tokenReppo := ctn.Get("persistence.token").(*persistence.TokenPersistence)
				return service.NewProductService(), nil
			},
	}...)
}
