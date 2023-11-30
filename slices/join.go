package slices

// Join joins the supplied slices into a single slice.
func Join[V any](s ...[]V) []V {
	totalItems := Reduce(s, 0, func(result int, next []V, _ int) int {
		return result + len(next)
	})
	arr := make([]V, 0, totalItems)
	for _, slice := range s {
		arr = append(arr, slice...)
	}
	return arr
}
