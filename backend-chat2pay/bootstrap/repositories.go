package bootstrap

import (
	"chat2pay/config/yaml"
	"chat2pay/internal/repositories"
	"github.com/jmoiron/sqlx"
	"github.com/sarulabs/di/v2"
)

func loadRepositories(builder *di.Builder, config *yaml.Config) {
	builder.Add([]di.Def{
		{
			Name: ProductRepositoriesName,
			Build: func(ctn di.Container) (interface{}, error) {
				return repositories.NewProductRepo(ctn.Get(DatabaseAdapter).(*sqlx.DB)), nil
			},
		},
	}...)
}
