package bootstrap

import (
	"chat2pay/internal/api/handlers"
	"chat2pay/internal/service"
	"github.com/sarulabs/di/v2"
)

func LoadHandler() *[]di.Def {
	return &[]di.Def{
		{
			Name: PaymentHandlerName,
			Build: func(ctn di.Container) (interface{}, error) {
				paymentService := ctn.Get(PaymentServiceName).(service.PaymentService)
				return handlers.NewPaymentHandler(paymentService), nil
			},
		},
		{
			Name: WebhookHandlerName,
			Build: func(ctn di.Container) (interface{}, error) {
				paymentService := ctn.Get(PaymentServiceName).(service.PaymentService)
				return handlers.NewWebhookHandler(paymentService), nil
			},
		},
	}
}
