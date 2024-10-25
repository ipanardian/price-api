package dtoV1

import "github.com/ipanardian/price-api/internal/model/frame"

type GetPriceRequest struct {
	Ids []string `validate:"required" query:"ids"`
}

type GetPriceResponse struct {
	ID    string            `json:"id"`
	Price frame.HermesPrice `json:"price"`
}
