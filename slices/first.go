package slices

// First gets the first element of the supplied slice, along with a boolean true
// retrieving the first value succeeded; false otherwise.
//
// False will only be returned if the slice is empty or null.
func First[T any](s []T) (T, bool) {
	var zero T
	if len(s) == 0 {
		return zero, false
	}
	return s[0], true
}

// FirstOrDefautl gets the first element of the supplied slice, if the slice is
// non-empty, otherwise, just returns the default value.
func FirstOrDefault[T any](s []T, value T) T {
	if len(s) == 0 {
		return value
	}
	return s[0]
}
