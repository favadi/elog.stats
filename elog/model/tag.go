package model

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
)

// Tags ...
type Tags map[string]string

// Value ...
func (t Tags) Value() (driver.Value, error) {
	bytes, err := json.Marshal(t)
	return string(bytes), err
}

// Scan ...
func (t *Tags) Scan(src interface{}) error {
	if src == nil {
		return nil
	}
	source, ok := src.([]byte)
	if !ok {
		return errors.New("tags: unable to convert data into binary")
	}
	if err := json.Unmarshal(source, t); err != nil {
		return err
	}
	return nil
}

// FromPbTags ...
func FromPbTags(mValues map[string]string) (Tags, error) {
	t := Tags{}
	for key, value := range mValues {
		t[key] = value
	}
	return t, nil
}

// ToPbTags ...
func ToPbTags(values Tags) (map[string]string, error) {
	t := map[string]string{}
	for key, value := range values {
		t[key] = value
	}
	return t, nil
}
