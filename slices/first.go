package slices

func First[T any](s []T) (T, bool) {
	var zero T
	if len(s) == 0 {
		return zero, false
	}
	return s[0], true
}

func FirstOrDefault[T any](s []T, value T) T {
	if len(s) == 0 {
		return value
	}
	return s[0]
}
