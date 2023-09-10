package slices

func Filter[T any](slice []T, predicate func(T, int) bool) []T {
	result := make([]T, 0, len(slice))
	for i, item := range slice {
		if predicate(item, i) {
			result = append(result, item)
		}
	}
	return result
}
