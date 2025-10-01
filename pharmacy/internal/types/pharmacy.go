package types

import (
	"encoding/json"
	"fmt"
)

type OpeningHours map[string][]string

func (oh *OpeningHours) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("invalid type for OpeningHours: %T", value)
	}
	return json.Unmarshal(b, oh)
}
