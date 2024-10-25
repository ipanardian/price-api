package types

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/shopspring/decimal"
)

type ArrayDecimal []decimal.Decimal

func NewArrayDecimal(values ...decimal.Decimal) ArrayDecimal {
	return values
}

func (t *ArrayDecimal) ToArray() []decimal.Decimal {
	return *t
}

func (t *ArrayDecimal) String() string {
	jsonStr, e := json.Marshal(t)
	if e != nil {
		return "[]"
	} else {
		return string(jsonStr)
	}
}

// Scan scan value into Jsonb, implements sql.Scanner interface
func (t *ArrayDecimal) Scan(value any) error {
	if value != nil {
		bytes, ok := value.([]byte)
		if !ok {
			return errors.New(fmt.Sprint("Failed to unmarshal JSONB value:", value))
		}

		var strList []string

		result := ArrayDecimal{}
		err := json.Unmarshal(bytes, &strList)
		for _, str := range strList {
			d, e := decimal.NewFromString(str)
			if e == nil {
				result = append(result, d)
			}
		}
		*t = result
		return err
	} else {
		*t = ArrayDecimal{}
		return nil
	}
}

// Value return json value, implement driver.Valuer interface
func (t ArrayDecimal) Value() (driver.Value, error) {
	return t.String(), nil
}
