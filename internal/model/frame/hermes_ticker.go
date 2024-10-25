package frame

import (
	"errors"
	"time"

	"github.com/shopspring/decimal"
)

type HermesResponse struct {
	Type      string          `json:"type"`
	PriceFeed HermesPriceFeed `json:"price_feed"`
	Status    string          `json:"status"`
	Error     string          `json:"error"`
}

type HermesPriceFeed struct {
	ID       string      `json:"id"`
	Price    HermesPrice `json:"price"`
	EmaPrice HermesPrice `json:"ema_price"`
	Vaa      string      `json:"vaa"`
}

type HermesPrice struct {
	Price       decimal.Decimal `json:"price"`
	Conf        decimal.Decimal `json:"conf"`
	Expo        int             `json:"expo"`
	PublishTime int64           `json:"publish_time"`
}

type HermesPairArg struct {
	Market string `json:"market"`
	Pair   string `json:"pair"`
	Pair0  string `json:"pair0"`
	Pair1  string `json:"pair1"`
}

func (p *HermesPriceFeed) GetPriceNoOlderThan(age int64) (price HermesPrice, err error) {
	currentTime := time.Now().Unix()

	if currentTime-p.Price.PublishTime > age {
		err = errors.New("price was expired")
		return
	}

	return p.Price, err
}
