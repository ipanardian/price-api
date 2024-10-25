//go:build wireinject
// +build wireinject

package injector

import (
	"github.com/google/wire"
	"github.com/ipanardian/price-api/internal/service"
)

func InitializeHermesService() (service.HermesService, error) {
	wire.Build(
		service.NewHermesService,
	)
	return nil, nil
}

func InitializeApiService() (service.ApiService, error) {
	wire.Build(
		service.NewApiService,
	)
	return nil, nil
}
