package types

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
)

type ArrayUint64 []uint64

func NewArrayUint64(values ...uint64) ArrayUint64 {
	return values
}

func (t *ArrayUint64) ToArray() []uint64 {
	return *t
}

func (t *ArrayUint64) String() string {
	jsonStr, e := json.Marshal(t)
	if e != nil {
		return "[]"
	} else {
		return string(jsonStr)
	}
}

// Scan scan value into Jsonb, implements sql.Scanner interface
func (t *ArrayUint64) Scan(value any) error {
	if value != nil {
		bytes, ok := value.([]byte)
		if !ok {
			return errors.New(fmt.Sprint("Failed to unmarshal JSONB value:", value))
		}

		result := ArrayUint64{}
		err := json.Unmarshal(bytes, &result)
		*t = result
		return err
	} else {
		*t = ArrayUint64{}
		return nil
	}
}

// Value return json value, implement driver.Valuer interface
func (t ArrayUint64) Value() (driver.Value, error) {
	return t.String(), nil
}
