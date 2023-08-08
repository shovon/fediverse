package slices

type Pusher[T any] []T

func (p *Pusher[T]) Push(v T) {
	*p = append(*p, v)
}
