package nullable

import (
	"encoding/json"
)

type Nilable[T any] struct {
	hasValue bool
	value    T
}

func Just[T any](v T) Nilable[T] {
	return Nilable[T]{
		hasValue: true,
		value:    v,
	}
}

func Nil[T any]() Nilable[T] {
	var t T
	return Nilable[T]{
		hasValue: false,
		value:    t,
	}
}

func (n Nilable[T]) HasValue() bool {
	return n.hasValue
}

func (n Nilable[T]) Value() (T, bool) {
	if !n.hasValue {
		return n.value, false
	}

	return n.value, true
}

func (n Nilable[T]) ValueOrDefault(d T) T {
	if !n.hasValue {
		return d
	}

	return n.value
}

func (n Nilable[T]) AssertValue() T {
	if !n.hasValue {
		panic("null dereference error")
	}
	return n.value
}

func (n Nilable[T]) MarshalJSON() ([]byte, error) {
	if !n.hasValue {
		return []byte("null"), nil
	}

	return json.Marshal(n.value)
}

func (n *Nilable[T]) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		n.hasValue = false
		var t T
		n.value = t
		return nil
	}

	var value T
	err := json.Unmarshal(data, &value)
	if err != nil {
		return err
	}

	n.hasValue = true
	n.value = value
	return nil
}

func Then[T any, V any](n Nilable[T], fn func(T) Nilable[V]) Nilable[V] {
	if !n.hasValue {
		return Nil[V]()
	}

	return fn(n.value)
}
