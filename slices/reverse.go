package slices

func Reverse[K any](slice []K) []K {
	result := make([]K, len(slice))
	for i := range slice {
		result[i] = slice[len(slice)-1-i]
	}
	return result
}
