//go:build wireinject
// +build wireinject

package injector

import (
	"github.com/google/wire"
	"github.com/ipanardian/price-api/internal/service"
)

func InitializeHermesService() (service.WsOracleService, error) {
	wire.Build(
		service.NewHermesService,
	)
	return nil, nil
}
