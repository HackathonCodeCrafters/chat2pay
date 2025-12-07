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
		{
			Name: MerchantRepositoryName,
			Build: func(ctn di.Container) (interface{}, error) {
				return repositories.NewMerchantRepo(ctn.Get(DatabaseAdapter).(*sqlx.DB)), nil
			},
		},
		{
			Name: CustomerRepositoryName,
			Build: func(ctn di.Container) (interface{}, error) {
				return repositories.NewCustomerRepo(ctn.Get(DatabaseAdapter).(*sqlx.DB)), nil
			},
		},
		{
			Name: MerchantUserRepositoryName,
			Build: func(ctn di.Container) (interface{}, error) {
				return repositories.NewMerchantUserRepo(ctn.Get(DatabaseAdapter).(*sqlx.DB)), nil
			},
		},
		{
			Name: OrderRepositoryName,
			Build: func(ctn di.Container) (interface{}, error) {
				return repositories.NewOrderRepository(ctn.Get(DatabaseAdapter).(*sqlx.DB)), nil
			},
		},
		{
			Name: ChatMessageRepositoryName,
			Build: func(ctn di.Container) (interface{}, error) {
				return repositories.NewChatMessageRepository(ctn.Get(DatabaseAdapter).(*sqlx.DB)), nil
			},
		},
	}
}
