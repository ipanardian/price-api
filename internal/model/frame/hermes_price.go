package frame

import (
	"encoding/json"

	"github.com/shopspring/decimal"
)

type PriceHermes struct {
	ID          string          `json:"id"`
	Price       decimal.Decimal `json:"price"`
	Conf        decimal.Decimal `json:"confidence"`
	Expo        int             `json:"expo"`
	PublishTime int64           `json:"publish_time"`
}

func (f PriceHermes) MarshalBinary() ([]byte, error) {
	return json.Marshal(f)
}

func (f *PriceHermes) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, f)
}
