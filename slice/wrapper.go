package slice

type Wrapper[T any] []T

func (w Wrapper[T]) ForAll(fn func(T) bool) bool {
	for _, item := range w {
		if !fn(item) {
			return false
		}
	}
	return true
}
