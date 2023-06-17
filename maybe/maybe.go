package maybe

import (
	"encoding/json"
	"errors"
)

type Maybe[T any] struct {
	value    T
	hasValue bool
}

var _ json.Marshaler = Maybe[int]{}
var _ json.Marshaler = Maybe[string]{}

func Just[T any](value T) Maybe[T] {
	return Maybe[T]{value, true}
}

func Nothing[T any]() Maybe[T] {
	return Maybe[T]{}
}

func (m Maybe[T]) Get() (T, error) {
	if !m.hasValue {
		var def T
		return def, errors.New("Nothing")
	}
	return m.value, nil
}

func (m Maybe[T]) MarshalJSON() ([]byte, error) {
	if !m.hasValue {
		return json.Marshal(nil)
	}
	return json.Marshal(m.value)
}

func (m *Maybe[T]) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		m.hasValue = false
		return nil
	}
	var value T
	err := json.Unmarshal(data, &value)
	if err != nil {
		return err
	}
	m.value = value
	m.hasValue = true
	return nil
}
