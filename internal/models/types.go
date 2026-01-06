package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
)

// StringList is a custom type for storing string slices in SQLite as JSON
type StringList []string

// Value converts StringList to database value (JSON string)
func (s StringList) Value() (driver.Value, error) {
	if s == nil {
		return "[]", nil
	}
	return json.Marshal(s)
}

// Scan converts database value to StringList
func (s *StringList) Scan(value interface{}) error {
	if value == nil {
		*s = []string{}
		return nil
	}

	var bytes []byte
	switch v := value.(type) {
	case []byte:
		bytes = v
	case string:
		bytes = []byte(v)
	default:
		return errors.New("failed to scan StringList")
	}

	return json.Unmarshal(bytes, s)
}
