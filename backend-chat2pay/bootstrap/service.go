package bootstrap

import (
	"chat2pay/internal/payment/xendit"
	"chat2pay/internal/repositories"
	"chat2pay/internal/service"
	"github.com/sarulabs/di/v2"
)

func LoadService() *[]di.Def {
	return &[]di.Def{
		{
			Name: PaymentServiceName,
			Build: func(ctn di.Container) (interface{}, error) {
				paymentRepo := ctn.Get(PaymentRepositoryName).(repositories.PaymentRepository)
				paymentLogRepo := ctn.Get(PaymentLogRepositoryName).(repositories.PaymentLogRepository)
				// OrderRepo set to nil - payment service handles this gracefully
				xenditClient := ctn.Get(XenditClientName).(*xendit.Client)
				return service.NewPaymentService(paymentRepo, paymentLogRepo, nil, xenditClient), nil
			},
		},
	}
}
