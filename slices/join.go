package slices

func Join[V any](s ...[]V) []V {
	totalItems := Reduce(s, 0, func(result int, next []V, _ int) int {
		return result + len(next)
	})
	arr := make([]V, 0, totalItems)
	index := 0
	for _, slice := range s {
		for _, element := range slice {
			arr[index] = element
			index++
		}
	}
	return arr
}
