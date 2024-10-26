package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/ansel1/merry/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/ipanardian/price-api/internal/cache"
	"github.com/ipanardian/price-api/internal/constant"
	dtoV1 "github.com/ipanardian/price-api/internal/dto/v1"
	"github.com/ipanardian/price-api/internal/logger"
	"github.com/ipanardian/price-api/internal/model/frame"
	"go.uber.org/zap"
)

type ApiServiceImpl struct {
}

func NewApiService() ApiService {
	return &ApiServiceImpl{}
}

func (s *ApiServiceImpl) GetPrice(c *fiber.Ctx, req *dtoV1.GetPriceRequest) (d []dtoV1.GetPriceResponse, statusNumber string, err error) {
	ids := req.Ids
	if len(ids) < 1 {
		statusNumber = constant.PriceUnavailable
		err = merry.Wrap(errors.New(StatusMap[statusNumber]), merry.WithMessage(StatusMap[statusNumber]), merry.WithHTTPCode(fiber.StatusNotFound))
		return
	}

	for _, id := range ids {
		func(id string) {
			key := fmt.Sprintf(constant.PriceCache, id)
			dataMem, ok := cache.MemGet(key)
			if ok {
				d = append(d, dataMem.(dtoV1.GetPriceResponse))
				return
			}

			f, prices := s.GetCachePrice(id)
			if !f {
				statusNumber = constant.PriceUnavailable
				err = merry.Wrap(errors.New(StatusMap[statusNumber]), merry.WithMessage(StatusMap[statusNumber]), merry.WithHTTPCode(fiber.StatusNotFound))
				return
			}

			pc := dtoV1.GetPriceResponse{
				ID: prices.ID,
				Price: frame.HermesPrice{
					Price:       prices.Price,
					Conf:        prices.Conf,
					Expo:        prices.Expo,
					PublishTime: prices.PublishTime,
				},
			}
			cache.MemSet(key, pc, time.Second*2)

			d = append(d, pc)
		}(id)
	}

	return
}

func (s *ApiServiceImpl) GetCachePrice(id string) (valid bool, price *frame.PriceHermes) {
	ctx, done := context.WithCancel(context.Background())
	defer done()
	qS := fmt.Sprintf("price:%s", id)

	val, e := cache.Get[[]byte](ctx, qS)
	if e != nil {
		logger.Log.Error("Error getting price from redis", zap.Error(e), zap.String("id", id))
		return false, nil
	}
	err := json.Unmarshal(val, &price)
	if err != nil {
		logger.Log.Error("Error parse price from redis", zap.Error(e), zap.String("id", id))
		return false, nil
	}

	return true, price
}
