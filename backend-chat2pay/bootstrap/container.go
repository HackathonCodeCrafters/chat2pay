package bootstrap

import (
	"chat2pay/config/yaml"
	"github.com/sarulabs/di/v2"
)

func NewContainer() (di.Container, error) {
	builder, err := di.NewBuilder()
	if err != nil {
		return di.Container{}, err
	}

	// Add Config
	if err := builder.Add(di.Def{
		Name: ConfigDefName,
		Build: func(ctn di.Container) (interface{}, error) {
			return yaml.NewConfig()
		},
	}); err != nil {
		return di.Container{}, err
	}

	// Adapter
	if err := builder.Add(*NewAdapter()...); err != nil {
		return di.Container{}, err
	}

	// Repository
	if err := builder.Add(*NewRepository()...); err != nil {
		return di.Container{}, err
	}

	// Repository
	if err := builder.Add(*LoadService()...); err != nil {
		return di.Container{}, err
	}

	return builder.Build(), nil
}
