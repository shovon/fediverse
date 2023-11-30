package slices

// Has returns true if the slice contains the value.
func Has[T comparable](s []T, value T) bool {
	for _, item := range s {
		if item == value {
			return true
		}
	}
	return false
}
