package slices

func Map[T any, V any](slice []T, fn func(T, int) V) []V {
	result := make([]V, len(slice))
	for i, item := range slice {
		result[i] = fn(item, i)
	}
	return result
}

// IgnoreIndex is just a helper function for converting unary functions to
// functions that can accept an index, even when they don't need to.
//
// This is especially useful when paired with the `Map` function.
//
// Usage:
//
//	Map(slice, IgnoreIndex(strings.TrimSpace))
//
// In the above example, `strings.TrimSpace“ is unary function which is
// incompatible with `Map`. IgnoreIndex is used to wrap such unary functions to
// be compatible with Map.
func IgnoreIndex[T, V any](f func(T) V) func(T, int) V {
	return func(v T, _ int) V {
		return f(v)
	}
}
