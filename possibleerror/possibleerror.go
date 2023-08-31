package possibleerror

import "encoding/json"

type PossibleError[T any] struct {
	value T
	err   error
}

var _ json.Marshaler = PossibleError[string]{}

func NotError[T any](value T) PossibleError[T] {
	return PossibleError[T]{value: value}
}

func Error[T any](err error) PossibleError[T] {
	return PossibleError[T]{err: err}
}

func New[T any](value T, err error) PossibleError[T] {
	return PossibleError[T]{value: value, err: err}
}

func Then[T any, V any](p PossibleError[T], fn func(T) PossibleError[V]) PossibleError[V] {
	if p.err != nil {
		return Error[V](p.err)
	}
	return fn(p.value)
}

// MapToThen is like `Then`, but returns a function that in-turn returns a
// PossibleError
func MapToThen[T any, V any](fn func(t T) V) func(t T) PossibleError[V] {
	return func(t T) PossibleError[V] {
		return NotError[V](fn(t))
	}
}

func (p PossibleError[T]) Value() (T, error) {
	return p.value, p.err
}

func (p PossibleError[T]) MarshalJSON() ([]byte, error) {
	if p.err != nil {
		return nil, p.err
	}
	return json.Marshal(p.value)
}
