package bootstrap

import (
	"chat2pay/internal/repositories"
	"github.com/jmoiron/sqlx"
	"github.com/sarulabs/di/v2"
)

func NewRepository() *[]di.Def {
	return &[]di.Def{
		{
			Name: PaymentRepositoryName,
			Build: func(ctn di.Container) (interface{}, error) {
				return repositories.NewPaymentRepo(ctn.Get(DatabaseAdapter).(*sqlx.DB)), nil
			},
		},
		{
			Name: PaymentLogRepositoryName,
			Build: func(ctn di.Container) (interface{}, error) {
				return repositories.NewPaymentLogRepo(ctn.Get(DatabaseAdapter).(*sqlx.DB)), nil
			},
		},
	}
}
