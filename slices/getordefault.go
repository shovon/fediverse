package slices

// GetOrDefault gets the element at the specified index, or returns the default
// value if the index is out of bounds.
func GetOrDefault[T any](s []T, index int, defaultValue T) T {
	if index < 0 || index >= len(s) {
		return defaultValue
	}
	return s[index]
}
