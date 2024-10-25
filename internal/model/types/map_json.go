package types

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"

	"google.golang.org/protobuf/types/known/structpb"
)

type MapJSON map[string]MutatedValue

func (t *MapJSON) String() string {
	jsonStr, e := json.Marshal(t)
	if e != nil {
		return ""
	} else {
		return string(jsonStr)
	}
}

// Scan scan value into Jsonb, implements sql.Scanner interface
func (t *MapJSON) Scan(value any) error {
	if value != nil {
		bytes, ok := value.([]byte)
		if !ok {
			return errors.New(fmt.Sprint("Failed to unmarshal JSONB value:", value))
		}

		result := MapJSON{}
		err := json.Unmarshal(bytes, &result)
		*t = result
		return err
	} else {
		*t = MapJSON{}
		return nil
	}
}

// Value return json value, implement driver.Valuer interface
func (t MapJSON) Value() (driver.Value, error) {
	if t == nil {
		return "{}", nil
	}
	return t.String(), nil
}

func (t *MapJSON) UnmarshalJSON(data []byte) error {
	var un map[string]any
	_ = json.Unmarshal(data, &un)

	*t = MapJSON{}
	for k, v := range un {
		(*t)[k] = MutatedValue(fmt.Sprintf("%v", v))
	}

	//fmt.Println(t)

	return nil
}

func (t *MapJSON) ToProtoBuff() (res *structpb.Struct, err error) {
	if t.String() == "" {
		return nil, nil
	}

	buff := structpb.Struct{}
	err = buff.UnmarshalJSON([]byte(t.String()))
	if err != nil {
		return nil, nil
	}

	return &buff, err
}

func (t *MapJSON) AddValue(key, value string) {
	if t == nil {
		*t = make(MapJSON)
	}

	(*t)[key] = MutatedValue(value)
}

func (t *MapJSON) GetValue(key string) MutatedValue {
	if t != nil {
		if val, ok := (*t)[key]; ok {
			return val
		}
	}
	return MutatedValue("")
}
