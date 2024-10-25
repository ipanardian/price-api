package handler

import "github.com/ipanardian/price-api/internal/service"

type Handler struct {
	Api *ApiHandler
}

func NewHandler(
	apiService service.ApiService,
) *Handler {
	return &Handler{
		Api: NewApiHandler(apiService),
	}
}
