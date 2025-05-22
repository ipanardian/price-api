package service

import (
	"github.com/gofiber/fiber/v2"
	"github.com/ipanardian/price-api/internal/constant"
	dtoV1 "github.com/ipanardian/price-api/internal/dto/v1"
)

var StatusMap = map[string]string{
	constant.InternalError:        "Internal server error.",
	constant.ApiMaintenance:       "API service is under maintenance.",
	constant.WebsocketMaintenance: "Websocket service is under maintenance.",
	constant.GrpcMaintenance:      "gRPC service is under maintenance.",
	constant.InvalidRequest:       "Invalid request.",
	constant.RestrictedUser:       "User is not allowed to access this resource.",
	constant.PriceNotFound:        "Price not found.",
	constant.PriceUnavailable:     "Price is unavailable.",
	constant.RequiredIdsParam:     "Required ids paramater.",
}

type ApiService interface {
	GetPrice(c *fiber.Ctx, req *dtoV1.GetPriceRequest) ([]dtoV1.GetPriceResponse, string, error)
}

type HermesService interface {
	Run()
	HealthCheck()
}
