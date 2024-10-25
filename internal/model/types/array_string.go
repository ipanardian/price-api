package types

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
)

type ArrayString []string

func NewArrayString(values ...string) ArrayString {
	return values
}

func (t *ArrayString) ToArray() []string {
	return *t
}

func (t *ArrayString) String() string {
	jsonStr, e := json.Marshal(t)
	if e != nil {
		return "[]"
	} else {
		return string(jsonStr)
	}
}

// Scan scan value into Jsonb, implements sql.Scanner interface
func (t *ArrayString) Scan(value any) error {
	if value != nil {
		bytes, ok := value.([]byte)
		if !ok {
			return errors.New(fmt.Sprint("Failed to unmarshal JSONB value:", value))
		}

		result := ArrayString{}
		err := json.Unmarshal(bytes, &result)
		*t = result
		return err
	} else {
		*t = ArrayString{}
		return nil
	}
}

// Value return json value, implement driver.Valuer interface
func (t ArrayString) Value() (driver.Value, error) {
	return t.String(), nil
}
