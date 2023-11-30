package iterable

// Iterable represents a type that can be iterated on
type Iterable[K any] interface {
	Iterate() <-chan K
}

// ToSlice takes an iterable, and converts it into a slice
func ToSlice[K any](i Iterable[K]) []K {
	result := []K{}

	for v := range i.Iterate() {
		result = append(result, v)
	}

	return result
}
