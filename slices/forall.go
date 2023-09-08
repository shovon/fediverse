package slices

func ForAll[T any](s []T, predicate func(T) bool) bool {
	for _, item := range s {
		if !predicate(item) {
			return false
		}
	}
	return true
}
