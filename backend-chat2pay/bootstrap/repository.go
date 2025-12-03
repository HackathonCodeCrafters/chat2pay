package bootstrap

import (
	"chat2pay/internal/repositories"
	"github.com/jmoiron/sqlx"
	"github.com/sarulabs/di/v2"
)

func NewRepository() *[]di.Def {
	return &[]di.Def{
		{
			Name: ProductRepositoryName,
			Build: func(ctn di.Container) (interface{}, error) {
				return repositories.NewProductRepo(ctn.Get(DatabaseAdapter).(*sqlx.DB)), nil
			},
		},
	}
}
