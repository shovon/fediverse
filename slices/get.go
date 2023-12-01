package slices

// Get gets the element at the specified index, or returns false if the index
// is out of bounds.
func Get[T any](s []T, index int) (T, bool) {
	var zero T
	if index < 0 || index >= len(s) {
		return zero, false
	}
	return s[index], true
}
