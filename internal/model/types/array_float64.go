package types

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
)

type ArrayFloat64 []float64

func NewArrayFloat64(values ...float64) ArrayFloat64 {
	return values
}

func (t *ArrayFloat64) ToArray() []float64 {
	return *t
}

func (t *ArrayFloat64) String() string {
	jsonStr, e := json.Marshal(t)
	if e != nil {
		return "[]"
	} else {
		return string(jsonStr)
	}
}

// Scan scan value into Jsonb, implements sql.Scanner interface
func (t *ArrayFloat64) Scan(value any) error {
	if value != nil {
		bytes, ok := value.([]byte)
		if !ok {
			return errors.New(fmt.Sprint("Failed to unmarshal JSONB value:", value))
		}

		result := ArrayFloat64{}
		err := json.Unmarshal(bytes, &result)
		*t = result
		return err
	} else {
		*t = ArrayFloat64{}
		return nil
	}
}

// Value return json value, implement driver.Valuer interface
func (t ArrayFloat64) Value() (driver.Value, error) {
	return t.String(), nil
}
