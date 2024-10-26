package handler

import (
	"fmt"
	"net/http"
	"time"

	"github.com/ansel1/merry/v2"
	"github.com/gofiber/fiber/v2"
	dtoV1 "github.com/ipanardian/price-api/internal/dto/v1"
	"github.com/ipanardian/price-api/internal/service"
)

type ApiHandler struct {
	ApiService service.ApiService
}

func NewApiHandler(apiService service.ApiService) *ApiHandler {
	return &ApiHandler{
		ApiService: apiService,
	}
}

func (h *ApiHandler) setResponseV1(c *fiber.Ctx, res dtoV1.ResponseWrapper) {
	code := http.StatusOK
	if res.Error != nil {
		res.Status = 0
		res.Data = nil
		code = http.StatusInternalServerError

		if res.StatusNumber == "" {
			res.StatusMessage = "Sorry, this is not working properly. We know about this mistake and are working to fix it."
		} else {
			res.StatusMessage = res.Error.Error()
		}

		if merry.HTTPCode(res.Error) != 0 {
			code = merry.HTTPCode(res.Error)
			res.StatusCode = fmt.Sprintf("%d", code)
		}
	} else {
		res.Status = 1
		res.StatusCode = fmt.Sprintf("%d", code)
		res.StatusNumber = res.StatusCode
		res.StatusMessage = "success"
	}

	res.Timestamp = time.Now().UTC().UnixMilli()
	c.Response().Header.SetStatusCode(code)
	c.JSON(res)
}

// GetPrice
// @summary Get price by price feed ids
// @description Get price by price feed ids
// @tags v1
// @produce json
// @router /v1/price [get]
// @security ClientIdAuth
// @security ClientSignatureAuth
// @param ids[] query []string true "Format: ?ids[]=a12...&ids[]=b4c..." collectionFormat(multi)
// @success 200 {object} dtoV1.ResponseWrapper{data=[]dtoV1.GetPriceResponse}
// @failure 500 {object} dtoV1.ResponseWrapper
func (h *ApiHandler) GetPrice(c *fiber.Ctx) (err error) {
	var req dtoV1.GetPriceRequest
	var res dtoV1.ResponseWrapper

	if err := c.QueryParser(&req); err != nil {
		return err
	}

	res.Data, res.StatusNumber, res.Error = h.ApiService.GetPrice(c, &req)
	h.setResponseV1(c, res)

	return
}
