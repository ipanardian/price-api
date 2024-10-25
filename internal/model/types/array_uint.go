package types

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
)

type ArrayUint []uint

func NewArrayUint(values ...uint) ArrayUint {
	return values
}

func (t *ArrayUint) ToArray() []uint {
	return *t
}

func (t *ArrayUint) String() string {
	jsonStr, e := json.Marshal(t)
	if e != nil {
		return "[]"
	} else {
		return string(jsonStr)
	}
}

// Scan scan value into Jsonb, implements sql.Scanner interface
func (t *ArrayUint) Scan(value any) error {
	if value != nil {
		bytes, ok := value.([]byte)
		if !ok {
			return errors.New(fmt.Sprint("Failed to unmarshal JSONB value:", value))
		}

		result := ArrayUint{}
		err := json.Unmarshal(bytes, &result)
		*t = result
		return err
	} else {
		*t = ArrayUint{}
		return nil
	}
}

// Value return json value, implement driver.Valuer interface
func (t ArrayUint) Value() (driver.Value, error) {
	return t.String(), nil
}
