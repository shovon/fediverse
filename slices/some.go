package slices

func Some[T any](s []T, fn func(T) bool) bool {
	for _, item := range s {
		if fn(item) {
			return true
		}
	}
	return false
}
