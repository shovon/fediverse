package slices

type Wrapper[T any] []T

func (w Wrapper[T]) ForAll(fn func(T) bool) bool {
	for _, item := range w {
		if !fn(item) {
			return false
		}
	}
	return true
}

func (w Wrapper[T]) Some(fn func(T) bool) bool {
	for _, item := range w {
		if fn(item) {
			return true
		}
	}
	return false
}

func (w Wrapper[T]) Map(fn func(T, int) T) Wrapper[T] {
	return Map(w, fn)
}
