package pointers

func ValueOrDefault[T any](value *T, def *T) *T {
	if value == nil {
		return def
	}
	return value
}
