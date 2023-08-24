package slices

func Map[T any, V any](slice []T, fn func(T) V) []V {
	result := make([]V, len(slice))
	for i, item := range slice {
		result[i] = fn(item)
	}
	return result
}
