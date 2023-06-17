package jrd

import (
	"encoding/json"
	"errors"
)

type NullableString struct {
	hasValue bool
	value    string
}

func JustString(v string) NullableString {
	return NullableString{
		hasValue: true,
		value:    v,
	}
}

func NullString() NullableString {
	return NullableString{
		hasValue: false,
		value:    "",
	}
}

func (n NullableString) HasValue() bool {
	return n.hasValue
}

func (n NullableString) Value() (string, error) {
	if !n.hasValue {
		return "", errors.New("no value")
	}

	return n.value, nil

}

func (n NullableString) MarshalJSON() ([]byte, error) {
	if !n.hasValue {
		return []byte("null"), nil
	}

	return json.Marshal(n.value)
}

func (n *NullableString) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		n.hasValue = false
		n.value = ""
		return nil
	}

	var value string
	err := json.Unmarshal(data, &value)
	if err != nil {
		return err
	}

	n.hasValue = true
	n.value = value
	return nil
}
