package possibleerror

import "encoding/json"

type PossibleError[T any] struct {
	value T
	err   error
}

var _ json.Marshaler = PossibleError[string]{}

func Value[T any](value T) PossibleError[T] {
	return PossibleError[T]{value: value}
}

func Error[T any](err error) PossibleError[T] {
	return PossibleError[T]{err: err}
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
